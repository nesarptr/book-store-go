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

func GetSingleBook(c *fiber.Ctx) error {
	bookId := c.Params("id")
	book := new(models.Book)
	db := config.GetDB()
	db.First(book, bookId)
	if book.Title == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid book id",
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "book successfully retrived",
		"data":    book,
	})
}

func GetCart(c *fiber.Ctx) error {
	return nil
}

func PostCart(c *fiber.Ctx) error {
	userId := c.Locals("userId").(float64)
	db := config.GetDB()
	cart := new(models.Cart)
	type cartBook struct {
		ID       string `json:"id"`
		Quantity int    `json:"quantity"`
	}
	type inputCart struct {
		Books []cartBook `json:"books"`
	}
	db.Where("user_id = ?", userId).First(cart)
	if cart.ID == 0 {
		cart.UserID = uint(userId)
		cart.TotalPrice = 0
		cart.Create(db)
	}
	InputCart := new(inputCart)
	if err := c.BodyParser(InputCart); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"data": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"cart": cart,
	})
}

func RemoveCart(c *fiber.Ctx) error {
	return nil
}
