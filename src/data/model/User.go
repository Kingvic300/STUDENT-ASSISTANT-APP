package model

import (
	    "go.mongodb.org/mongo-driver/bson/primitive"
		"Student-Assistant-App/src/data/enums"
)


type User struct {
    ID       	primitive.ObjectID 		`bson:"_id,omitempty"   json:"id"`
    Name     	string            		`bson:"name"            json:"name"`
    Email    	string            		`bson:"email"           json:"email"`
    Password 	string            		`bson:"password"        json:"-"`
    Role     	enums.Role          	`bson:"role"            json:"role"`
}

func (req*User) SetUser(ID primitive.ObjectID){
    req.ID = ID
}
func (req*User) GetUser() primitive.ObjectID{
    return req.ID
}
func (req*User) SetName(Name string){
    req.Name = Name
}
func (req*User) GetName() string{
    return req.Name
}
func (req*User) SetEmail(Email string){
    req.Email = Email
}
func (req*User) GetEmail() string{
    return req.Email
}
func (req*User) SetPassword(Password string){
    req.Password = Password
}
func (req*User) GetPassword() string{
    return req.Password
}
func (req *User) SetRole(Role enums.Role){
    req.Role = Role
}
func (req *User) GetRole() enums.Role{
    return req.Role
}
