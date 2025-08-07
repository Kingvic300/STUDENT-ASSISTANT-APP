package model

import (
	    "go.mongodb.org/mongo-driver/bson/primitive"
		"Student-Assistant-App/src/data/enums"
)


type User struct {
    ID       	primitive.ObjectID 		`bson:"_id,omitempty" json:"id"`
    Name     	string            		`bson:"name" json:"name"`
    Email    	string            		`bson:"email" json:"email"`
    Password 	string            		`bson:"password" json:"-"`
    Role     	enums.Role          	`bson:"role" json:"role"`
}