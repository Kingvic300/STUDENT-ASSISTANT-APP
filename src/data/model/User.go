package model

import "Student-Assistant-App/src/data/enums"


type User struct {
	Id   		uint 			`json:"id"`
	Name 		string 			`json:"name"`
	Email 		string 			`json:"email"`
	Password	string 			`json:"-"`
	Role  		enums.Role		`json:"Role"`
}