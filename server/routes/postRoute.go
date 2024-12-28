package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/johnj09/golang-mongodb-crud-project/controllers"
	"github.com/johnj09/golang-mongodb-crud-project/middleware"
)

func PostRoute(app *fiber.App) {
	app.Post("/api/post", middleware.Authenticate, controllers.CreatePost)
	app.Get("/api/post/:pid", controllers.GetPost)
	app.Get("/api/posts/:idx", controllers.GetNextTenPosts)
	app.Patch("/api/post/:pid", middleware.Authenticate, controllers.EditPost)
	app.Delete("/api/post/:pid", middleware.Authenticate, controllers.DeletePost)
}