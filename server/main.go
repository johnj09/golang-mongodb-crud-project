package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/johnj09/golang-mongodb-crud-project/configs"
	"github.com/johnj09/golang-mongodb-crud-project/routes"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file: ", err)
	}

	app := fiber.New()

	configs.ConnectDB()

	routes.UserRoute(app)
	routes.PostRoute(app)
	routes.CommentRoute(app)
	
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	log.Fatal(app.Listen("0.0.0.0:" + port))
}