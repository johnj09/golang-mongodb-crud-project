package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Comment struct {
	ID     		primitive.ObjectID  `json:"id,omitempty" bson:"_id,omitempty"`
	PostID 		primitive.ObjectID 	`json:"post_id" bson:"post_id"`
	UserID 		primitive.ObjectID 	`json:"user_id" bson:"user_id"`
	Content 	string 				`json:"content" bson:"content"`
	CreatedAt 	int64				`json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt	int64				`json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

type CommentEditRequest struct {
	NewContent string `json:"content"`
}