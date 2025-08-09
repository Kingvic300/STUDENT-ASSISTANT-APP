package response

import "Student-Assistant-App/src/data/model"

type Response[T any] struct {
	Data    T    `json:"data"`
	Success bool `json:"success"`
}

type CreateUserResponse struct {
	Message string      `json:"message"`
	User    *model.User `json:"user"`
	Token   string      `json:"token,omitempty"`
}

func (req *CreateUserResponse) SetMessage(Message string) {
	req.Message = Message
}
func (req *CreateUserResponse) GetMessage() string {
	return req.Message
}
func (req *CreateUserResponse) SetUser(User *model.User) {
	req.User = User
}
func (req *CreateUserResponse) GetUser() *model.User {
	return req.User
}
func (req *CreateUserResponse) SetToken(token string) {
	req.Token = token
}
func (req *CreateUserResponse) GetToken() string {
	return req.Token
}

type LoginResponse struct {
	Message string      `json:"message"`
	User    *model.User `json:"user"`
	Token   string      `json:"token"`
}

func (r *LoginResponse) SetMessage(message string) {
	r.Message = message
}
func (r *LoginResponse) GetMessage() string {
	return r.Message
}
func (r *LoginResponse) SetUser(user *model.User) {
	r.User = user
}
func (r *LoginResponse) GetUser() *model.User {
	return r.User
}
func (r *LoginResponse) SetToken(token string) {
	r.Token = token
}
func (r *LoginResponse) GetToken() string {
	return r.Token
}

type DeleteUserResponse struct {
	Message string `json:"message"`
}

func (r *DeleteUserResponse) SetMessage(message string) {
	r.Message = message
}
func (r *DeleteUserResponse) GetMessage() string {
	return r.Message
}
