package controllers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/nesarptr/book-store-go/config"
	"github.com/nesarptr/book-store-go/models"
	"github.com/nesarptr/book-store-go/utils"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(c *fiber.Ctx) error {
	user := new(models.User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	errors := utils.ValidateStruct(*user)

	if errors != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(errors)
	}

	if err := user.Create(config.GetDB()); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"id":    user.ID,
		"email": user.Email,
	})
}

func SignIn(c *fiber.Ctx) error {
	type loginInput struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=6"`
	}

	input := new(loginInput)

	if err := c.BodyParser(input); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	errors := utils.ValidateStruct(*input)

	if errors != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(errors)
	}

	user := new(models.User)

	config.GetDB().Where("email = ?", input.Email).First(user)

	if user.Email == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "user does not exist with this email",
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "password is not valid",
		})
	}

	claims := jwt.MapClaims{
		"id":    user.ID,
		"email": user.Email,
		"exp":   time.Now().Add(time.Minute * 15).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secret, err := config.GetEnv("TOKEN_SECRET")

	if err != nil {
		return fiber.ErrInternalServerError
	}

	t, err := token.SignedString([]byte(secret))

	if err != nil {
		return fiber.ErrInternalServerError
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": "Success login", "token": t})
}
