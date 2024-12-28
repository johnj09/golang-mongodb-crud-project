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
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()
	
	userId, err := primitive.ObjectIDFromHex(c.Locals("userId").(string))
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	// Retrieve request body
	post := new(models.Post)
	if err := c.BodyParser(post); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	// Check for valid title and content
	if post.Title == "" || post.Content == "" {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": "Post title and content must not be empty."}})
	}

	// Create a post post id
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
	result, err := postCollection.InsertOne(ctx, newPost)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusCreated).JSON(responses.UserResponse{Status: http.StatusCreated, Message: "success", Data: &fiber.Map{"data": result}})
}

func GetPost(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	postId := c.Params("id")
	var post models.Post
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(postId)

	err := postCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&post)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusOK).JSON(responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": post}})
}

func GetNextTenPosts(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	idx, _ := strconv.Atoi(c.Params("idx"))
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}}).SetLimit(10).SetSkip(int64(idx * 10))

	cursor, err := postCollection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}
	defer cursor.Close(ctx)

	var posts []models.Post
	if err := cursor.All(ctx, &posts); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	return c.Status(http.StatusOK).JSON(responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": posts}})
}

func UpdatePostContent(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	postId,_ := primitive.ObjectIDFromHex(c.Params("id"))
	updatedContent := new(models.PostUpdateRequest)
	if err := c.BodyParser(updatedContent); err != nil {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: &fiber.Map{"data": err.Error()}})
	}

	filter := bson.M{"_id": postId}
	update := bson.M{
		"$set": bson.M{
			"content": updatedContent.UpdatedContent,
			"updated_at": time.Now().Unix(),
		},
	}

	result, err :=  postCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": "Failed to update post."}})
	}

	if result.MatchedCount == 0 {
		return c.Status(http.StatusBadRequest).JSON(responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: &fiber.Map{"data": "No post found with given post ID."}})
	}

	return c.Status(http.StatusOK).JSON(responses.UserResponse{Status: http.StatusOK, Message: "success", Data: &fiber.Map{"data": result}})

}