package responses

import (
	"github.com/gofiber/fiber/v2"
)

type UserResponse struct {
	Status  int			`json:"status"`
	Message string     	`json:"message"`
	Data    *fiber.Map 	`json:"data"`
}

// func newUserResponse(c *fiber.Ctx, statusCode int, msg string, data string) error {
// 	return c.Status(statusCode).JSON(UserResponse{Status: statusCode, Message: msg, Data: &fiber.Map{"data": data}})
// }