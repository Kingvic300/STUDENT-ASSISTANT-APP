package request
import (
		"Student-Assistant-App/src/data/enums"
)

type CreateUserRequest struct {
	Name     	string            		`bson:"name" json:"name"`
    Email    	string            		`bson:"email" json:"email"`
    Password 	string            		`bson:"password" json:"-"`
    Role     	enums.Role          	`bson:"role" json:"role"`
}