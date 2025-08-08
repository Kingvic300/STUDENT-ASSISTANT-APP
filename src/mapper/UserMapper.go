package mapper

import (
	"Student-Assistant-App/src/data/model"
	"Student-Assistant-App/src/dtos/request"
	"Student-Assistant-App/src/utils"
)

func MapToUser(req *request.CreateUserRequest) (*model.User, error) {
	email, err := utils.EmailVerification(req.Email)
	if err != nil {
		return nil, err
	}

	return &model.User{
		Password: req.Password,
		Email:    email,
		Role:     req.Role,
	}, nil
}
