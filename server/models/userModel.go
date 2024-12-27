package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID       	primitive.ObjectID 	`json:"id,omitempty" bson:"_id,omitempty"`
	Username 	string             	`json:"username"`
	Email    	string             	`json:"email"`
	Password 	string             	`json:"password"`
	CreatedAt 	time.Time			`json:"created_at"`
}