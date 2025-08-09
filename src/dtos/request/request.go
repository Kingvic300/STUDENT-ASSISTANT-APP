package request

import (
	"Student-Assistant-App/src/data/enums"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CreateUserRequest struct {
	Name     string     `bson:"name"     json:"name"`
	Email    string     `bson:"email"    json:"email"`
	Password string     `bson:"password" json:"password"`
	Role     enums.Role `bson:"role"     json:"role"`
}

func (req *CreateUserRequest) SetName(Name string) {
	req.Name = Name
}
func (req *CreateUserRequest) GetName() string {
	return req.Name
}
func (req *CreateUserRequest) SetEmail(Email string) {
	req.Email = Email
}
func (req *CreateUserRequest) GetEmail() string {
	return req.Email
}
func (req *CreateUserRequest) SetPassword(Password string) {
	req.Password = Password
}
func (req *CreateUserRequest) GetPassword() string {
	return req.Password
}
func (req *CreateUserRequest) SetRole(Role enums.Role) {
	req.Role = Role
}
func (req *CreateUserRequest) GetRole() enums.Role {
	return req.Role
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (req *LoginRequest) SetEmail(email string) {
	req.Email = email
}
func (req *LoginRequest) GetEmail() string {
	return req.Email
}
func (req *LoginRequest) SetPassword(password string) {
	req.Password = password
}
func (req *LoginRequest) GetPassword() string {
	return req.Password
}

type SendOTPRequest struct {
	Email   string `json:"email" binding:"required"`
	Purpose string `json:"purpose" binding:"required"` // "signup", "login", "password_reset"
}

func (req *SendOTPRequest) SetEmail(email string) {
	req.Email = email
}
func (req *SendOTPRequest) GetEmail() string {
	return req.Email
}
func (req *SendOTPRequest) SetPurpose(purpose string) {
	req.Purpose = purpose
}
func (req *SendOTPRequest) GetPurpose() string {
	return req.Purpose
}

type VerifyOTPRequest struct {
	Email   string `json:"email" binding:"required"`
	Code    string `json:"code" binding:"required"`
	Purpose string `json:"purpose" binding:"required"`
}

func (req *VerifyOTPRequest) SetEmail(email string) {
	req.Email = email
}
func (req *VerifyOTPRequest) GetEmail() string {
	return req.Email
}
func (req *VerifyOTPRequest) SetCode(code string) {
	req.Code = code
}
func (req *VerifyOTPRequest) GetCode() string {
	return req.Code
}
func (req *VerifyOTPRequest) SetPurpose(purpose string) {
	req.Purpose = purpose
}
func (req *VerifyOTPRequest) GetPurpose() string {
	return req.Purpose
}

type SignupWithOTPRequest struct {
	CreateUserRequest
	OTPCode string `json:"otp_code" binding:"required"`
}

func (req *SignupWithOTPRequest) SetOTPCode(code string) {
	req.OTPCode = code
}
func (req *SignupWithOTPRequest) GetOTPCode() string {
	return req.OTPCode
}

type LoginWithOTPRequest struct {
	Email   string `json:"email" binding:"required"`
	OTPCode string `json:"otp_code" binding:"required"`
}

func (req *LoginWithOTPRequest) SetEmail(email string) {
	req.Email = email
}
func (req *LoginWithOTPRequest) GetEmail() string {
	return req.Email
}
func (req *LoginWithOTPRequest) SetOTPCode(code string) {
	req.OTPCode = code
}
func (req *LoginWithOTPRequest) GetOTPCode() string {
	return req.OTPCode
}

type UpdateUserRequest struct {
	Name  string     `json:"name"`
	Email string     `json:"email"`
	Role  enums.Role `json:"role"`
}

func (req *UpdateUserRequest) SetName(name string) {
	req.Name = name
}
func (req *UpdateUserRequest) GetName() string {
	return req.Name
}
func (req *UpdateUserRequest) SetEmail(email string) {
	req.Email = email
}
func (req *UpdateUserRequest) GetEmail() string {
	return req.Email
}
func (req *UpdateUserRequest) SetRole(role enums.Role) {
	req.Role = role
}
func (req *UpdateUserRequest) GetRole() enums.Role {
	return req.Role
}

type DeleteUserRequest struct {
	Id primitive.ObjectID `json:"id" bson:"_id"`
}

func (req *DeleteUserRequest) SetId(ID primitive.ObjectID) {
	req.Id = ID
}
func (req *DeleteUserRequest) GetId() primitive.ObjectID {
	return req.Id
}