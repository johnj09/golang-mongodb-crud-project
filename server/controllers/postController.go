package controllers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/johnj09/golang-mongodb-crud-project/configs"
	"github.com/johnj09/golang-mongodb-crud-project/models"
	"github.com/johnj09/golang-mongodb-crud-project/responses"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var postCollection = configs.GetCollection("posts")

func CreatePost(c *fiber.Ctx) error {
	userId, err := primitive.ObjectIDFromHex(c.Locals("userId").(string))
	if err != nil {
		return responses.NewUserResponse(c, http.StatusInternalServerError, responses.Error, err.Error())
	}

	// Retrieve request body
	post := new(models.Post)
	if err := c.BodyParser(post); err != nil {
		return responses.NewUserResponse(c, http.StatusBadRequest, responses.Error, err.Error())
	}

	// Check for valid title and content
	if post.Title == "" || post.Content == "" {
		return responses.NewUserResponse(c, http.StatusBadRequest, responses.Error, "Post title and content must not be empty.")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	// Create a unique post id
	var postId primitive.ObjectID
	for duplicate := true; duplicate; {
		postId = primitive.NewObjectID()
		filter := bson.M{"_id": postId}
		existingId := userCollection.FindOne(ctx, filter)
		if existingId.Err() == mongo.ErrNoDocuments {
			duplicate = false
		}
	}

	newPost := models.Post{
		ID: postId,
		UserID: userId,
		Title: post.Title,
		Content: post.Content,
		NumOfLikes: 0,
		NumOfComments: 0,
		CreatedAt: time.Now().Unix(),
		UpdatedAt: time.Now().Unix(),
	}

	// Insert new post into database
	_, err = postCollection.InsertOne(ctx, newPost)
	if err != nil {
		return responses.NewUserResponse(c, http.StatusInternalServerError, responses.Error, err.Error())
	}

	return responses.NewUserResponse(c, http.StatusCreated, responses.Success, "Post Created.")
}

func GetPost(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	postId := c.Params("pid")
	var post models.Post
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(postId)

	err := postCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&post)
	if err != nil {
		return responses.NewUserResponse(c, http.StatusInternalServerError, responses.Error, err.Error())
	}

	return responses.NewUserResponse(c, http.StatusOK, responses.Success, post)
}

func GetNextTenPosts(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	idx, _ := strconv.Atoi(c.Params("idx"))
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}).SetLimit(10).SetSkip(int64(idx * 10))

	cursor, err := postCollection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return responses.NewUserResponse(c, http.StatusInternalServerError, responses.Error, err.Error())
	}
	defer cursor.Close(ctx)

	var posts []models.Post
	if err := cursor.All(ctx, &posts); err != nil {
		return responses.NewUserResponse(c, http.StatusInternalServerError, responses.Error, err.Error())
	}

	return responses.NewUserResponse(c, http.StatusOK, responses.Success, posts)
}

func EditPost(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	// Convert postId to ObjectID
	postId,_ := primitive.ObjectIDFromHex(c.Params("pid"))
	editRequest := new(models.PostEditRequest)
	if err := c.BodyParser(editRequest); err != nil {
		return responses.NewUserResponse(c, http.StatusBadRequest, responses.Error, err.Error())
	}

	// Check user rights
	userId, _ := primitive.ObjectIDFromHex(c.Locals("userId").(string))
	var post models.Post
	err := postCollection.FindOne(ctx, bson.M{"_id": postId}).Decode(&post)
	if err == mongo.ErrNoDocuments {
		return responses.NewUserResponse(c, http.StatusBadRequest, responses.Error, err.Error())
	} else if err != nil {
		return responses.NewUserResponse(c, http.StatusInternalServerError, responses.Error, err.Error())
	}

	if post.UserID != userId {
		return responses.NewUserResponse(c, http.StatusUnauthorized, responses.Error, "Unauthorized.")
	}

	// Update post record
	filter := bson.M{"_id": postId}
	update := bson.M{
		"$set": bson.M{
			"content": editRequest.NewContent,
			"updated_at": time.Now().Unix(),
		},
	}

	result, err :=  postCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return responses.NewUserResponse(c, http.StatusInternalServerError, responses.Error, "Failed to update post.")
	}

	if result.MatchedCount == 0 {
		return responses.NewUserResponse(c, http.StatusInternalServerError, responses.Error, "No post found with given post ID.")
	}

	return responses.NewUserResponse(c, http.StatusOK, responses.Success, result)
}

func DeletePost(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	// Convert postId to ObjectID
	postId,_ := primitive.ObjectIDFromHex(c.Params("pid"))

	// Check user rights
	userId, _ := primitive.ObjectIDFromHex(c.Locals("userId").(string))
	var post models.Post
	err := postCollection.FindOne(ctx, bson.M{"_id": postId}).Decode(&post)
	if err == mongo.ErrNoDocuments {
		return responses.NewUserResponse(c, http.StatusBadRequest, responses.Error, err.Error())
	} else if err != nil {
		return responses.NewUserResponse(c, http.StatusInternalServerError, responses.Error, err.Error())
	}

	if post.UserID != userId {
		return responses.NewUserResponse(c, http.StatusUnauthorized, responses.Error, "Unauthorized.")
	}

	// Delete post record
	_, err = postCollection.DeleteOne(ctx, bson.M{"_id": postId})
	if err != nil {
		return responses.NewUserResponse(c, http.StatusInternalServerError, responses.Error, err.Error())
	}

	// Delete all comments from that post
	_, err = commentCollection.DeleteMany(ctx, bson.M{"post_id": postId})
	if err != nil {
		return responses.NewUserResponse(c, http.StatusInternalServerError, responses.Error, err.Error())
	}

	return responses.NewUserResponse(c, http.StatusOK, responses.Success, "Successfully deleted post.")
}