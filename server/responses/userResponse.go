package responses

import (
	"github.com/gofiber/fiber/v2"
)

const (
	Success = 0
	Error = 1
)

type UserResponse[T any] struct {
	Message 	string     	`json:"message"`
	Data    	T 			`json:"data"`
}

func NewUserResponse[T any](c *fiber.Ctx, statusCode int, message int, data T) error {
	var strMessage string
	if message == Success {
		strMessage = "success"
	} else if message == Error {
		strMessage = "error"
	}

	return c.Status(statusCode).JSON(UserResponse[T]{Message: strMessage, Data: data})
}