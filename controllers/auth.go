package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nesarptr/book-store-go/config"
	"github.com/nesarptr/book-store-go/models"
	"github.com/nesarptr/book-store-go/utils"
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

	return c.Status(fiber.StatusCreated).JSON(user)
}
