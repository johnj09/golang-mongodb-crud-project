package controllers

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/johnj09/golang-mongodb-crud-project/configs"
	"github.com/johnj09/golang-mongodb-crud-project/models"
	"github.com/johnj09/golang-mongodb-crud-project/responses"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection = configs.GetCollection("users")
var jwtSecret = []byte(configs.EnvSecretKey())

func CreateUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	// Retrieve request body
	user := new(models.User)
	if err := c.BodyParser(user); err != nil {
		return responses.NewUserResponse(c, http.StatusBadRequest, responses.Error, err.Error())
	}

	if user.Email == "" || user.Username == "" || user.Password == "" {
		return responses.NewUserResponse(c, http.StatusBadRequest, responses.Error, "Fields may not be empty.")
	}

	// Check for duplicate users
	filter := bson.M{"$or": []bson.M{
		{"email": user.Email},
		{"username": user.Username},
	}}
	existingUser := userCollection.FindOne(ctx, filter)
	if err := existingUser.Err(); err == nil {
		return responses.NewUserResponse(c, http.StatusConflict, responses.Error, "Email or username already exists.")
	} else if err != mongo.ErrNoDocuments {
		return responses.NewUserResponse(c, http.StatusInternalServerError, responses.Error, err.Error())
	}
	
	// Create Hashed Password
	salt := make([]byte, 16)
	rand.Read(salt)
	passwordWithSalt := append([]byte(user.Password)[:], salt[:]...)
	hashedPassword, err := bcrypt.GenerateFromPassword(passwordWithSalt, bcrypt.DefaultCost)
	if err != nil {
		return responses.NewUserResponse(c, http.StatusInternalServerError, responses.Error, "Internal Server Error.")
	}

	// Create a unique user id
	var userId primitive.ObjectID
	for duplicate := true; duplicate; {
		userId = primitive.NewObjectID()
		filter = bson.M{"_id": userId}
		existingId := userCollection.FindOne(ctx, filter)
		if existingId.Err() == mongo.ErrNoDocuments {
			duplicate = false
		}
	}

	newUser := models.User{
		ID: userId,
		Username: user.Username,
		Email: user.Email,
		Password: hex.EncodeToString(hashedPassword),
		Salt: hex.EncodeToString(salt),
		CreatedAt: time.Now().Unix(),
	}

	// Insert new user into database
	result, err := userCollection.InsertOne(ctx, newUser)
	if err != nil {
		return responses.NewUserResponse(c, http.StatusInternalServerError, responses.Error, err.Error())
	}

	return responses.NewUserResponse(c, http.StatusCreated, responses.Success, result)
}

func LoginUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()

	loginData := new(models.LoginRequest)
	if err := c.BodyParser(loginData); err != nil {
		return responses.NewUserResponse(c, http.StatusBadRequest, responses.Error, err.Error())
	}
	
	// Validate login info
	var user models.User
	err := userCollection.FindOne(ctx, bson.M{"email": loginData.Email}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return responses.NewUserResponse(c, http.StatusUnauthorized, responses.Error, "Invalid Email.")
	} else if err != nil {
		return responses.NewUserResponse(c, http.StatusInternalServerError, responses.Error, err.Error())
	}

	userPassword, _ := hex.DecodeString(user.Password)
	salt, _ := hex.DecodeString(user.Salt) 
	if bcrypt.CompareHashAndPassword(userPassword, append([]byte(loginData.Password)[:], salt[:]...)) != nil {
		return responses.NewUserResponse(c, http.StatusUnauthorized, responses.Error, "Invalid password.")
	}
	
	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": user.ID,
		"username": user.Username,
		"exp": time.Now().Add(30 * time.Minute).Unix(),
	})
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return responses.NewUserResponse(c, http.StatusInternalServerError, responses.Error, "Failed to create token.")
	}

	// Create cookie
	c.Cookie(&fiber.Cookie{
		Name: configs.AuthCookie,
		Value: tokenString,
		Expires: time.Now().Add(24 * time.Hour),
		Secure: false,
		HTTPOnly: true,
	})

	result := &fiber.Map{
		"userId": user.ID,
		"username": user.Username,
	}

	return responses.NewUserResponse(c, http.StatusOK, responses.Success, result)
}

func LogoutUser(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name: configs.AuthCookie,
		Value: "",
		Expires: time.Now().Add(-1 * time.Hour),
	})
	return responses.NewUserResponse(c, http.StatusOK, responses.Success, "Logout Success.")
}

func GetUserProfile(c *fiber.Ctx) error {
	tokenString := c.Cookies(configs.AuthCookie)
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(configs.EnvSecretKey()), nil
	})
	claims := token.Claims.(jwt.MapClaims)
	return responses.NewUserResponse(c, http.StatusOK, responses.Success, claims)
}

func GetUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	userId := c.Params("id")
	var user models.User
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(userId)

	err := userCollection.FindOne(ctx, bson.M{"_id": objId}).Decode(&user)
	if err != nil {
		return responses.NewUserResponse(c, http.StatusInternalServerError, responses.Error, err.Error())
	}

	return responses.NewUserResponse(c, http.StatusOK, responses.Success, user)
}

func DeleteUser(c *fiber.Ctx) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	userId := c.Params("id")
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(userId)

	result, err := userCollection.DeleteOne(ctx, bson.M{"_id": objId})
	if err != nil {
		return responses.NewUserResponse(c, http.StatusInternalServerError, responses.Error, err.Error())
	}

	if result.DeletedCount < 1 {
		return responses.NewUserResponse(c, http.StatusNotFound, responses.Error, "User with specified ID not found.")
	}

	return responses.NewUserResponse(c, http.StatusOK, responses.Success, "User successfully deleted.")
}