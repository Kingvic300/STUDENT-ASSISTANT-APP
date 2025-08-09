package service

import (
	"Student-Assistant-App/src/data/model"
	"Student-Assistant-App/src/dtos/request"
	"Student-Assistant-App/src/dtos/response"
	"Student-Assistant-App/src/utils"
	"context"
	"errors"
)

type AuthService interface {
	Login(ctx context.Context, request *request.LoginRequest) (*response.LoginResponse, error)
	GenerateTokenForUser(user *model.User) (string, error)
}

type AuthServiceImpl struct {
	userService UserService
}

func NewAuthService(userService UserService) AuthService {
	return &AuthServiceImpl{
		userService: userService,
	}
}

func (auth *AuthServiceImpl) Login(ctx context.Context, request *request.LoginRequest) (*response.LoginResponse, error) {
	if request.Email == "" {
		return nil, errors.New("email is required")
	}
	if request.Password == "" {
		return nil, errors.New("password is required")
	}

	user, err := auth.userService.GetUserByEmail(ctx, request.Email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("invalid email or password")
	}

	if !utils.CheckPassword(request.Password, user.Password) {
		return nil, errors.New("invalid email or password")
	}

	token, err := utils.GenerateJWT(user.ID.Hex(), user.Email, user.Role)
	if err != nil {
		return nil, err
	}

	return &response.LoginResponse{
		Message: "Login successful",
		User:    user,
		Token:   token,
	}, nil
}

func (auth *AuthServiceImpl) GenerateTokenForUser(user *model.User) (string, error) {
	return utils.GenerateJWT(user.ID.Hex(), user.Email, user.Role)
}