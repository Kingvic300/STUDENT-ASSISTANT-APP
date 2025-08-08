package service

import (
	"Student-Assistant-App/src/dtos/request"
	"Student-Assistant-App/src/dtos/response"
	"context"

)

type UserService interface {
	CreateUser(ctx context.Context, request *request.CreateUserRequest) (*response.CreateUserResponse, error)
	DeleteUser(ctx context.Context, request )
}