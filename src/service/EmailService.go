package service

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"

	"gopkg.in/gomail.v2"
)

type EmailService interface {
	SendOTP(email, otp, purpose string) error
	SendWelcomeEmail(email, name string) error
}

type EmailServiceImpl struct {
	host     string
	port     int
	username string
	password string
	from     string
	dialer   *gomail.Dialer
	once     sync.Once
}

func NewEmailService() (EmailService, error) {
	host := os.Getenv("EMAIL_HOST")
	portStr := os.Getenv("EMAIL_PORT")
	username := os.Getenv("EMAIL_USERNAME")
	password := os.Getenv("EMAIL_PASSWORD")
	from := os.Getenv("EMAIL_FROM")

	if host == "" || portStr == "" || username == "" || password == "" || from == "" {
		return nil, fmt.Errorf("email configuration missing")
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("invalid email port: %v", err)
	}

	dialer := gomail.NewDialer(host, port, username, password)


	return &EmailServiceImpl{
		host:     host,
		port:     port,
		username: username,
		password: password,
		from:     from,
		dialer:   dialer,
	}, nil
}

func (e *EmailServiceImpl) SendOTP(email, otp, purpose string) error {
	subject := "Your OTP Code"
	var body string

	switch purpose {
	case "signup":
		body = fmt.Sprintf(`
			<html>
			<body>
				<h2>Welcome to Student Assistant App!</h2>
				<p>Thank you for signing up. Please use the following OTP to verify your email address:</p>
				<div style="background-color: #f0f0f0; padding: 20px; text-align: center; font-size: 24px; font-weight: bold; color: #333; border-radius: 5px; margin: 20px 0;">
					%s
				</div>
				<p>This OTP will expire in 10 minutes.</p>
				<p>If you didn't request this, please ignore this email.</p>
			</body>
			</html>
		`, otp)
	case "login":
		body = fmt.Sprintf(`
			<html>
			<body>
				<h2>Login Verification</h2>
				<p>Please use the following OTP to complete your login:</p>
				<div style="background-color: #f0f0f0; padding: 20px; text-align: center; font-size: 24px; font-weight: bold; color: #333; border-radius: 5px; margin: 20px 0;">
					%s
				</div>
				<p>This OTP will expire in 10 minutes.</p>
				<p>If you didn't request this, please ignore this email and secure your account.</p>
			</body>
			</html>
		`, otp)
	case "password_reset":
		body = fmt.Sprintf(`
			<html>
			<body>
				<h2>Password Reset</h2>
				<p>You requested to reset your password. Please use the following OTP:</p>
				<div style="background-color: #f0f0f0; padding: 20px; text-align: center; font-size: 24px; font-weight: bold; color: #333; border-radius: 5px; margin: 20px 0;">
					%s
				</div>
				<p>This OTP will expire in 10 minutes.</p>
				<p>If you didn't request this, please ignore this email.</p>
			</body>
			</html>
		`, otp)
	default:
		body = fmt.Sprintf(`
			<html>
			<body>
				<h2>Your OTP Code</h2>
				<p>Please use the following OTP:</p>
				<div style="background-color: #f0f0f0; padding: 20px; text-align: center; font-size: 24px; font-weight: bold; color: #333; border-radius: 5px; margin: 20px 0;">
					%s
				</div>
				<p>This OTP will expire in 10 minutes.</p>
			</body>
			</html>
		`, otp)
	}

	return e.sendEmail(email, subject, body)
}

func (e *EmailServiceImpl) SendWelcomeEmail(email, name string) error {
	subject := "Welcome to Student Assistant App!"
	body := fmt.Sprintf(`
		<html>
		<body>
			<h2>Welcome %s!</h2>
			<p>Your account has been successfully verified and created.</p>
			<p>You can now enjoy all the features of Student Assistant App.</p>
			<p>Thank you for joining us!</p>
		</body>
		</html>
	`, name)

	return e.sendEmail(email, subject, body)
}

func (e *EmailServiceImpl) sendEmail(to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", e.from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	if err := e.dialer.DialAndSend(m); err != nil {
		log.Printf("failed to send email to %s: %v", to, err)
		return fmt.Errorf("failed to send email: %w", err)
	}

	log.Printf("email sent successfully to %s with subject %s", to, subject)
	return nil
}
