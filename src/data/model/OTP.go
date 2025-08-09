package model

import (
	"time"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OTP struct {
	ID        	primitive.ObjectID 		`bson:"_id,omitempty"json:"id"`
	Email     	string             		`bson:"email"json:"email"`
	Code      	string             		`bson:"code"json:"code"`
	Purpose   	string             		`bson:"purpose"json:"purpose"`
	ExpiresAt 	time.Time          		`bson:"expires_at"json:"expires_at"`
	Used      	bool               		`bson:"used"json:"used"`
	CreatedAt 	time.Time         		`bson:"created_at"json:"created_at"`
}

func (otp *OTP) IsExpired() bool {
	return time.Now().After(otp.ExpiresAt)
}

func (otp *OTP) IsValid() bool {
	return !otp.Used && !otp.IsExpired()
}