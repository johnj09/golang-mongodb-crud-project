package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/johnj09/golang-mongodb-crud-project/configs"
	"github.com/johnj09/golang-mongodb-crud-project/responses"
)

func Authenticate(c *fiber.Ctx) error {
	// Get the cookie from the request
	tokenString := c.Cookies(configs.AuthCookie)
	if tokenString == "" {
		return c.Status(http.StatusUnauthorized).JSON(responses.UserResponse{Status: http.StatusUnauthorized, Message: "error", Data: &fiber.Map{"data": "Unauthorized"}})
	}

	// Decode and validate token
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(configs.EnvSecretKey()), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		// Check the expiration
		exp, ok := claims["exp"].(float64)
		if !ok || float64(time.Now().Unix()) > exp {
			return c.Status(http.StatusUnauthorized).JSON(responses.UserResponse{Status: http.StatusUnauthorized, Message: "error", Data: &fiber.Map{"data": "Token expired."}})
		}

		userId, ok := claims["userId"].(string)
		if !ok {
			return c.Status(http.StatusUnauthorized).JSON(responses.UserResponse{Status: http.StatusUnauthorized, Message: "error", Data: &fiber.Map{"data": "Invalid token payload."}})
		}

		c.Locals("userId", userId)
	} else {
		return c.Status(http.StatusUnauthorized).JSON(responses.UserResponse{Status: http.StatusUnauthorized, Message: "error", Data: &fiber.Map{"data": "Unauthorized."}})
	}

	return c.Next()	
}