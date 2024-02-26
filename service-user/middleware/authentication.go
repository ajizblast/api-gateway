package middleware

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"service-user/config"
	"service-user/helpers"
	"service-user/model"

	"github.com/gofiber/fiber/v2"
)

func Authentication(c *fiber.Ctx) error {
	access_token := c.Get("access_token")

	if len(access_token) == 0 {
		return c.Status(401).SendString("Invalid token: Access token missing")
	}

	checkToken, err := helpers.VerifyToken(access_token)

	if err != nil {
		return c.Status(401).SendString("Invalid token: Failed to verify token")
	}

	email := checkToken["email"].(string)

	db := config.GetDB()

	var user model.User
	result := db.Where("email = ?", email).First(&user)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return c.Status(401).SendString("Invalid token: User not found")
	} else if result.Error != nil {
		fmt.Println(err, "Error fetching user from database")
		return c.Status(500).SendString("Internal server error")
	}

	// Set user data in context for future use
	c.Locals("user", user)

	// Continue processing if user is found
	return c.Next()
}
