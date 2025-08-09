package service

import (
	"Student-Assistant-App/src/data/model"
	"Student-Assistant-App/src/data/repository"
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"time"
)

type OTPService interface {
	GenerateAndSendOTP(ctx context.Context, email, purpose string) error
	VerifyOTP(ctx context.Context, email, code, purpose string) error
	ResendOTP(ctx context.Context, email, purpose string) error
}

type OTPServiceImpl struct {
	otpRepository repository.OTPRepository
	emailService  EmailService
}

func NewOTPService(otpRepo repository.OTPRepository, emailService EmailService) OTPService {
	return &OTPServiceImpl{
		otpRepository: otpRepo,
		emailService:  emailService,
	}
}

func (s *OTPServiceImpl) GenerateAndSendOTP(ctx context.Context, email, purpose string) error {
	otpCode, err := s.generateOTP()
	if err != nil {
		return err
	}

	otp := &model.OTP{
		Email:     email,
		Code:      otpCode,
		Purpose:   purpose,
		ExpiresAt: time.Now().Add(2 * time.Minute), 
		Used:      false,
	}

	_, err = s.otpRepository.Save(ctx, otp)
	if err != nil {
		return err
	}

	return s.emailService.SendOTP(email, otpCode, purpose)
}

func (s *OTPServiceImpl) VerifyOTP(ctx context.Context, email, code, purpose string) error {
	otp, err := s.otpRepository.FindByEmailAndCode(ctx, email, code)
	if err != nil {
		return err
	}

	if otp == nil {
		return errors.New("invalid or expired OTP")
	}

	if otp.Purpose != purpose {
		return errors.New("OTP purpose mismatch")
	}

	if !otp.IsValid() {
		return errors.New("invalid or expired OTP")
	}

	return s.otpRepository.MarkAsUsed(ctx, otp.ID)
}

func (s *OTPServiceImpl) ResendOTP(ctx context.Context, email, purpose string) error {
	latestOTP, err := s.otpRepository.FindLatestByEmailAndPurpose(ctx, email, purpose)
	if err != nil {
		return err
	}

	if latestOTP != nil && time.Since(latestOTP.CreatedAt) < time.Minute {
		return errors.New("please wait before requesting a new OTP")
	}

	return s.GenerateAndSendOTP(ctx, email, purpose)
}

func (s *OTPServiceImpl) generateOTP() (string, error) {
	max := big.NewInt(999999)
	min := big.NewInt(100000)
	
	n, err := rand.Int(rand.Reader, max.Sub(max, min).Add(max, big.NewInt(1)))
	if err != nil {
		return "", err
	}
	
	return fmt.Sprintf("%06d", n.Add(n, min).Int64()), nil
}