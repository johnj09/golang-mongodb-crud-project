package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/johnj09/golang-mongodb-crud-project/controllers"
	"github.com/johnj09/golang-mongodb-crud-project/middleware"
)

func CommentRoute(app *fiber.App) {
	app.Post("/api/comment/:pid", middleware.Authenticate, controllers.CreateComment)
	app.Get("/api/comment/:pid", controllers.GetAllComments)
	app.Patch("/api/comment/:cid", middleware.Authenticate, controllers.EditComment)
	app.Delete("/api/comment/:cid", middleware.Authenticate, controllers.DeleteComment)
}