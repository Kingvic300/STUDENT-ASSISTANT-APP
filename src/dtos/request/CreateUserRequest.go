package request
import (
		"Student-Assistant-App/src/data/enums"
)

type CreateUserRequest struct {
	Name     	string            		`bson:"name"        json:"name"`
    Email    	string            		`bson:"email"       json:"email"`
    Password 	string            		`bson:"password"    json:"-"`
    Role     	enums.Role          	`bson:"role"        json:"role"`
}
func (req *CreateUserRequest) SetName(Name string){
    req.Name = Name
}
func (req *CreateUserRequest) GetName() string{
    return req.Name
}
func (req *CreateUserRequest) SetEmail(Email string){
    req.Email = Email
}
func(req *CreateUserRequest) GetEmail() string {
    return req.Email
}
func(req *CreateUserRequest) SetPassword(Password string){
    req.Password = Password
}
func (req *CreateUserRequest) GetPassword() string{
    return req.Password
}
func (req *CreateUserRequest) SetRole(Role enums.Role){
    req.Role = Role
}
func (req *CreateUserRequest) GetRole() enums.Role{
    return req.Role
}