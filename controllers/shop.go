package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nesarptr/book-store-go/config"
	"github.com/nesarptr/book-store-go/models"
)

func GetAllBooks(c *fiber.Ctx) error {
	db := config.GetDB()
	books := new([]models.Book)
	db.Find(books)
	if len(*books) > 0 {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "all books retrived successfully",
			"data":    books,
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "no book found",
	})
}
