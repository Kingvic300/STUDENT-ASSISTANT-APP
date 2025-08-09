package service

import (
	"Student-Assistant-App/src/data/model"
	"Student-Assistant-App/src/data/repository"
	"Student-Assistant-App/src/dtos/request"
	"Student-Assistant-App/src/dtos/response"
	"Student-Assistant-App/src/mapper"
	"Student-Assistant-App/src/utils"
	"context"
	"errors"
)

type UserService interface {
	CreateUser(ctx context.Context, request *request.CreateUserRequest) (*response.CreateUserResponse, error)
	GetUserByID(ctx context.Context, id string) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	GetAllUsers(ctx context.Context) ([]*model.User, error)
	UpdateUser(ctx context.Context, id string, request *request.UpdateUserRequest) (*model.User, error)
	DeleteUser(ctx context.Context, id string) error
}

type UserServiceImpl struct {
	userRepository repository.UserRepository
}

func NewUserServiceImpl(userRepo repository.UserRepository) UserService {
	return &UserServiceImpl{
		userRepository: userRepo,
	}
}

func (userService *UserServiceImpl) CreateUser(ctx context.Context, request *request.CreateUserRequest) (*response.CreateUserResponse, error) {
	if request.Email == "" {
		return nil, errors.New("email is required")
	}
	if request.Name == "" {
		return nil, errors.New("name is required")
	}
	if request.Password == "" {
		return nil, errors.New("password is required")
	}

	existingUser, err := userService.userRepository.FindByEmail(ctx, request.Email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, errors.New("user already exists with this email")
	}

	user, err := mapper.MapToUser(request)
	if err != nil {
		return nil, err
	}

	user.Name = request.Name

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		return nil, err
	}
	user.Password = hashedPassword

	savedUser, err := userService.userRepository.Save(ctx, user)
	if err != nil {
		return nil, err
	}

	token, err := utils.GenerateJWT(savedUser.ID.Hex(), savedUser.Email, savedUser.Role)
	if err != nil {
		return nil, err
	}

	response := &response.CreateUserResponse{
		User:    savedUser,
		Message: "User created successfully",
		Token:   token,
	}
	return response, nil
}

func (userService *UserServiceImpl) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	return userService.userRepository.FindByID(ctx, id)
}

func (userService *UserServiceImpl) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	return userService.userRepository.FindByEmail(ctx, email)
}

func (userService *UserServiceImpl) GetAllUsers(ctx context.Context) ([]*model.User, error) {
	return userService.userRepository.FindAll(ctx)
}

func (userService *UserServiceImpl) UpdateUser(ctx context.Context, id string, request *request.UpdateUserRequest) (*model.User, error) {
	existingUser, err := userService.userRepository.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existingUser == nil {
		return nil, errors.New("user not found")
	}

	if request.Name != "" {
		existingUser.Name = request.Name
	}
	if request.Email != "" {
		validEmail, err := utils.EmailVerification(request.Email)
		if err != nil {
			return nil, err
		}
		userWithEmail, err := userService.userRepository.FindByEmail(ctx, validEmail)
		if err != nil {
			return nil, err
		}
		if userWithEmail != nil && userWithEmail.ID != existingUser.ID {
			return nil, errors.New("email already taken by another user")
		}
		existingUser.Email = validEmail
	}
	if request.Role.IsValid() {
		existingUser.Role = request.Role
	}

	return userService.userRepository.Save(ctx, existingUser)
}

func (userService *UserServiceImpl) DeleteUser(ctx context.Context, id string) error {
	existingUser, err := userService.userRepository.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if existingUser == nil {
		return errors.New("user not found")
	}

	return userService.userRepository.DeleteByID(ctx, id)
}
