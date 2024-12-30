package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Post struct {
	ID 				primitive.ObjectID 	`json:"id,omitempty" bson:"_id,omitempty"`
	UserID 			primitive.ObjectID	`json:"user_id" bson:"user_id"`
	Title 			string				`json:"title" bson:"title"`
	Content			string				`json:"content" bson:"content"`
	NumOfLikes		uint32				`json:"num_of_likes" bson:"num_of_likes"`
	NumOfComments	uint32				`json:"num_of_comments" bson:"num_of_comments"`
	CreatedAt		int64				`json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt		int64				`json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

type PostEditRequest struct {
	NewContent	string	`json:"new_content"`
}

type GetPostResult struct {
	ID 				primitive.ObjectID 	`json:"id"`
	UserID			primitive.ObjectID  `json:"user_id"`
	Username		string				`json:"username"`
	Title 			string				`json:"title"`
	Content			string				`json:"content"`
	NumOfLikes		uint32				`json:"num_of_likes"`
	NumOfComments	uint32				`json:"num_of_comments"`
	CreatedAt		int64				`json:"created_at"`
	UpdatedAt		int64				`json:"updated_at"`
}