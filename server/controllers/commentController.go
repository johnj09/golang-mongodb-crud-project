package controllers

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/johnj09/golang-mongodb-crud-project/configs"
	"github.com/johnj09/golang-mongodb-crud-project/models"
	"github.com/johnj09/golang-mongodb-crud-project/responses"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var commentCollection = configs.GetCollection("comments")

func CreateComment(c *fiber.Ctx) error {
	userId, err := primitive.ObjectIDFromHex(c.Locals("userId").(string))
	if err != nil {
		return responses.NewUserResponse(c, http.StatusInternalServerError, responses.Error, err.Error())
	}

	postId, err := primitive.ObjectIDFromHex(c.Params("pid"))
	if err != nil {
		return responses.NewUserResponse(c, http.StatusInternalServerError, responses.Error, err.Error())
	}

	// Retrieve request body
	comment := new(models.Comment)
	if err := c.BodyParser(comment); err != nil {
		return responses.NewUserResponse(c, http.StatusBadRequest, responses.Error, err.Error())
	}

	if strings.TrimSpace(comment.Content) == "" {
		return responses.NewUserResponse(c, http.StatusBadRequest, responses.Error, "Comment can not be empty.")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	// Create a unique comment id
	var commentId primitive.ObjectID
	for duplicate := true; duplicate; {
		commentId = primitive.NewObjectID()
		filter := bson.M{"_id": commentId}
		existingId := userCollection.FindOne(ctx, filter)
		if existingId.Err() == mongo.ErrNoDocuments {
			duplicate = false
		}
	}

	newComment := models.Comment{
		ID: commentId,
		PostID: postId,
		UserID: userId,
		Content: comment.Content,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}

	// Insert new comment into database
	result, err := commentCollection.InsertOne(ctx, newComment)
	if err != nil {
		return responses.NewUserResponse(c, http.StatusInternalServerError, responses.Error, err.Error())
	}

	// Increment num of comments by 1 in corresponding post document
	increment := bson.M{
		"$inc": bson.M{"num_of_comments": 1},
	}
	_, err = postCollection.UpdateOne(ctx, bson.M{"_id": postId}, increment)
	if err != nil {
		return responses.NewUserResponse(c, http.StatusInternalServerError, responses.Error, err.Error())
	}

	return responses.NewUserResponse(c, http.StatusCreated, responses.Success, result)
}

func GetAllComments(c *fiber.Ctx) error {
	postId, err := primitive.ObjectIDFromHex(c.Params("pid"))
	if err != nil {
		return responses.NewUserResponse(c, http.StatusInternalServerError, responses.Error, err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	cursor, err := commentCollection.Find(ctx, bson.M{"post_id": postId})
	if err != nil {
		return responses.NewUserResponse(c, http.StatusInternalServerError, responses.Error, err.Error())
	}
	defer cursor.Close(ctx)

	var comments []models.Comment
	if err := cursor.All(ctx, &comments); err != nil {
		return responses.NewUserResponse(c, http.StatusInternalServerError, responses.Error, err.Error())
	}

	return responses.NewUserResponse(c, http.StatusOK, responses.Success, comments)
}

func EditComment(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	// Convert commentId to ObjectID
	commentId,_ := primitive.ObjectIDFromHex(c.Params("cid"))
	editRequest := new(models.CommentEditRequest)
	if err := c.BodyParser(editRequest); err != nil {
		return responses.NewUserResponse(c, http.StatusBadRequest, responses.Error, err.Error())
	}

	// Check user rights
	userId, _ := primitive.ObjectIDFromHex(c.Locals("userId").(string))
	var comment models.Comment
	err := commentCollection.FindOne(ctx, bson.M{"_id": commentId}).Decode(&comment)
	if err == mongo.ErrNoDocuments {
		return responses.NewUserResponse(c, http.StatusBadRequest, responses.Error, err.Error())
	} else if err != nil {
		return responses.NewUserResponse(c, http.StatusInternalServerError, responses.Error, err.Error())
	}

	if comment.UserID != userId {
		return responses.NewUserResponse(c, http.StatusUnauthorized, responses.Error, "Unauthorized.")
	}

	// Update comment record
	filter := bson.M{"_id": commentId}
	update := bson.M{
		"$set": bson.M{
			"content": editRequest.NewContent,
			"updated_at": time.Now().Unix(),
		},
	}

	result, err :=  commentCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return responses.NewUserResponse(c, http.StatusInternalServerError, responses.Error, "Failed to update comment.")
	}

	if result.MatchedCount == 0 {
		return responses.NewUserResponse(c, http.StatusBadRequest, responses.Error, "No post found with given comment ID.")
	}

	return responses.NewUserResponse(c, http.StatusOK, responses.Success, result)
}

func DeleteComment(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	// Convert commentId to ObjectID
	commentId,_ := primitive.ObjectIDFromHex(c.Params("cid"))

	// Check user rights
	userId, _ := primitive.ObjectIDFromHex(c.Locals("userId").(string))
	var comment models.Comment
	err := commentCollection.FindOne(ctx, bson.M{"_id": commentId}).Decode(&comment)
	if err == mongo.ErrNoDocuments {
		return responses.NewUserResponse(c, http.StatusBadRequest, responses.Error, err.Error())
	} else if err != nil {
		return responses.NewUserResponse(c, http.StatusInternalServerError, responses.Error, err.Error())
	}

	if comment.UserID != userId {
		return responses.NewUserResponse(c, http.StatusUnauthorized, responses.Error, "Unauthorized.")
	}

	// Delete comment record
	_, err = commentCollection.DeleteOne(ctx, bson.M{"_id": commentId})
	if err != nil {
		return responses.NewUserResponse(c, http.StatusInternalServerError, responses.Error, "Unauthorized.")
	}

	return responses.NewUserResponse(c, http.StatusOK, responses.Success, "Successfully deleted comment.")
}