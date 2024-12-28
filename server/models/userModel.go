package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID       		primitive.ObjectID 	`json:"id,omitempty" bson:"_id,omitempty"`
	Username 		string             	`json:"username" bson:"username"`
	Email    		string             	`json:"email" bson:"email"`
	Password 		string             	`json:"password" bson:"password"`
	Salt			string				`json:"salt,omitempty" bson:"salt,omitempty"`
	CreatedAt 		int64 /*Unix Time*/	`json:"created_at,omitempty" bson:"created_at,omitempty"`
}

type LoginRequest struct {
	Email		string	`json:"email"`
	Password	string	`json:"password"`
}