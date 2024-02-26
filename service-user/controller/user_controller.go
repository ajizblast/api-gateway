package controller

import (
	"errors"
	"gorm.io/gorm"
	"service-user/helpers"
	"service-user/model"
	"strings"

	"service-user/config"
	
	"github.com/gofiber/fiber/v2"
)

type WebResponse struct {
	Code   int
	Status string
	Data   interface{}
}

func Register(c *fiber.Ctx) error {

	db := config.GetDB()
	var requestBody model.User

	c.BodyParser(&requestBody)

	hashedPassword := helpers.HashPassword([]byte(requestBody.Password))

	requestBody.Password = hashedPassword

	result := db.Create(&requestBody)

	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "uni_users_email") {
			return c.Status(400).JSON(WebResponse{
				Code:   400,
				Status: "BAD_REQUEST",
				Data:   "User already exists",
			})
		} else {
			return c.Status(500).JSON(WebResponse{
				Code:   500,
				Status: "BAD_REQUEST",
				Data:   result.Error.Error(),
			})
		}
	}

	access_token := helpers.SignToken(requestBody.Email)

	return c.JSON(struct {
		Code        int
		Status      string
		AccessToken string
		Data        interface{}
	}{
		Code:        200,
		Status:      "OK",
		AccessToken: access_token,
		Data:        requestBody,
	})
}

func Login(c *fiber.Ctx) error {
	db := config.GetDB()

	var requestBody model.User

	c.BodyParser(&requestBody)

	var user model.User
	result := db.Where("email = ?", requestBody.Email).First(&user)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return c.Status(401).JSON(WebResponse{
			Code:   401,
			Status: "BAD_REQUEST",
			Data:   "User not found",
		})
	} else if result.Error != nil {
		return c.Status(500).JSON(WebResponse{ // Unexpected database error
			Code:   500,
			Status: "BAD_REQUEST",
			Data:   result.Error.Error(),
		})
	}

	checkPassword := helpers.ComparePassword([]byte(user.Password), []byte(requestBody.Password))
	if !checkPassword {
		return c.JSON(WebResponse{
			Code:   401,
			Status: "BAD_REQUEST",
			Data:   errors.New("invalid password").Error(),
		})
	}

	access_token := helpers.SignToken(requestBody.Email)

	return c.JSON(struct {
		Code        int
		Status      string
		AccessToken string
		Data        interface{}
	}{
		Code:        200,
		Status:      "OK",
		AccessToken: access_token,
		Data:        user,
	})
}

func Auth(c *fiber.Ctx) error {
	return c.JSON("OK")
}
