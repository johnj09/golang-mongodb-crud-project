package responses

import (
	"github.com/gofiber/fiber/v2"
)

const (
	Success = 0
	Error = 1
)

type UserResponse struct {
	Status  	int			`json:"status"`
	Outcome 	string     	`json:"outcome"`
	Data    	*fiber.Map 	`json:"data"`
}

func NewUserResponse[T any](c *fiber.Ctx, statusCode int, outcome int, data T) error {
	var strOutcome string
	if outcome == Success {
		strOutcome = "success"
	} else if outcome == Error {
		strOutcome = "error"
	}
	return c.Status(statusCode).JSON(UserResponse{Status: statusCode, Outcome: strOutcome, Data: &fiber.Map{"data": data}})
}