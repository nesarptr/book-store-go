package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nesarptr/book-store-go/config"
	"github.com/nesarptr/book-store-go/models"
	"github.com/nesarptr/book-store-go/utils"
)

func CreateBook(c *fiber.Ctx) error {
	book := new(models.Book)
	if err := c.BodyParser(book); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"message": err.Error(),
		})
	}

	userId := c.Locals("userId").(float64)

	book.UserID = uint(userId)

	errors := utils.ValidateStruct(*book)

	if errors != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(errors)
	}

	db := config.GetDB()

	if err := book.Create(db); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	user := new(models.User)
	db.First(user, book.UserID)
	user.Books = append(user.Books, *book)
	db.Save(user)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data":    book,
		"message": "book created successfully",
	})
}
