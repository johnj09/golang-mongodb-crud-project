package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/johnj09/golang-mongodb-crud-project/controllers"
	"github.com/johnj09/golang-mongodb-crud-project/middleware"
)

func UserRoute(app *fiber.App) {
	app.Post("/api/register", controllers.CreateUser)
	app.Post("/api/login", controllers.LoginUser)
	app.Post("/api/logout", controllers.LogoutUser)
	app.Get("/api/profile", middleware.Authenticate, controllers.GetUserProfile)
	app.Get("/api/user/:id", controllers.GetUser)
	app.Delete("/api/user/:id", controllers.DeleteUser)
}