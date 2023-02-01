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

func GetBooks(c *fiber.Ctx) error {
	userId := c.Locals("userId").(float64)
	db := config.GetDB()
	user := new(models.User)
	db.Model(&models.User{}).Preload("Books").First(user, userId)
	if len(user.Books) > 0 {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "books retrived successfully",
			"data":    user.Books,
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "user does not own any book",
	})
}

func GetBook(c *fiber.Ctx) error {
	userId := c.Locals("userId").(float64)
	bookId := c.Params("id")
	db := config.GetDB()
	book := new(models.Book)
	db.First(book, bookId)
	if book.Title == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid book id",
		})
	}
	if book.UserID == uint(userId) {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "book successfully retrived",
			"data":    book,
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "this book does not belong to user",
	})
}

func UpdateBook(c *fiber.Ctx) error {
	userId := c.Locals("userId").(float64)
	bookId := c.Params("id")
	db := config.GetDB()
	book := new(models.Book)
	db.First(book, bookId)
	if book.Title == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid book id",
		})
	}
	if book.UserID == uint(userId) {
		if err := c.BodyParser(book); err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"data": err.Error(),
			})
		}

		errors := utils.ValidateStruct(*book)

		if errors != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(errors)
		}
		db.Save(book)

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "book successfully updated",
			"data":    book,
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "this book does not belong to user",
	})
}

func DeleteBook(c *fiber.Ctx) error {
	userId := c.Locals("userId").(float64)
	bookId := c.Params("id")
	db := config.GetDB()
	book := new(models.Book)
	db.First(book, bookId)
	if book.Title == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid book id",
		})
	}
	if book.UserID == uint(userId) {
		db.Unscoped().Where("book_id = ?", bookId).Delete(models.CartItem{})
		db.Unscoped().Delete(models.Book{}, bookId)
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "book successfully deleted",
			"data":    book,
		})
	}
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"message": "this book does not belong to user",
	})
}
