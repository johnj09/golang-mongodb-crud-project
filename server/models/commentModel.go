package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Comment struct {
	ID     		string             	`json:"id,omitempty" bson:"_id,omitempty"`
	PostID 		primitive.ObjectID 	`json:"post_id" bson:"post_id"`
	UserID 		primitive.ObjectID 	`json:"user_id" bson:"user_id"`
	Content 	string 				`json:"content" bson:"content"`
	CreatedAt 	time.Time 			`json:"created_at" bson:"created_at"`
	UpdatedAt	time.Time			`json:"updated_at" bson:"updated_at"`
}