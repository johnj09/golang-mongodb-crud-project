package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Post struct {
	ID 			primitive.ObjectID 	`json:"id,omitempty" bson:"_id,omitempty"`
	Author 		primitive.ObjectID 	`json:"author"`
	Title 		string				`json:"title"`
	Content		string				`json:"content"`
	CreatedAt	time.Time			`json:"created_at"`
	UpdatedAt	time.Time			`json:"updated_at"`
}