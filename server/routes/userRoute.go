package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/johnj09/golang-mongodb-crud-project/controllers"
)

func UserRoute(app *fiber.App) {
	app.Post("/api/register", controllers.CreateUser)
	app.Post("/api/login", controllers.LoginUser)
	app.Get("/api/user/:id", controllers.GetUser)
	app.Delete("/api/user/:id", controllers.DeleteUser)
}