package controllers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/nesarptr/book-store-go/config"
	"github.com/nesarptr/book-store-go/models"
	"github.com/nesarptr/book-store-go/utils"
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
	userId := c.Locals("userId").(float64)
	db := config.GetDB()
	cart := new(models.Cart)
	db.Preload("Books.Book").Where("user_id = ?", userId).First(cart)
	if cart.ID != 0 {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "cart successfully retrived",
			"data":    cart,
		})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "cart successfully retrived",
		"data":    "cart is empty",
	})
}

func PostCart(c *fiber.Ctx) error {
	userId := c.Locals("userId").(float64)
	db := config.GetDB()
	type cartBook struct {
		ID       string `json:"id"`
		Quantity int    `json:"quantity"`
	}
	type inputCart struct {
		Books []cartBook `json:"books"`
	}
	InputCart := new(inputCart)
	if err := c.BodyParser(InputCart); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"data": err.Error(),
		})
	}
	if len(InputCart.Books) <= 0 {
		return fiber.ErrBadRequest
	}
	cart := new(models.Cart)
	db.Where("user_id = ?", userId).First(cart)
	if cart.ID != 0 {
		db.Unscoped().Where("cart_id = ?", cart.ID).Delete(&models.CartItem{})
		db.Unscoped().Delete(cart)
	}
	cart.UserID = uint(userId)
	cart.TotalPrice = 0
	cart.Create(db)

	cartItems := make([]models.CartItem, 0)

	cartBooks := map[string]cartBook{}

	for _, book := range InputCart.Books {
		cartBooks[book.ID] = book
	}

	for _, book := range cartBooks {
		cartItem := new(models.CartItem)
		bookId, err := strconv.ParseUint(book.ID, 10, 64)
		if err != nil {
			return fiber.ErrInternalServerError
		}
		cartItem.CartID = cart.ID
		cartItem.BookID = uint(bookId)
		cartItem.Quantity = book.Quantity
		errors := utils.ValidateStruct(cartItem)
		if errors != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(errors)
		}
		cartItems = append(cartItems, *cartItem)
	}

	if db.Create(&cartItems).Error != nil {
		return fiber.ErrBadRequest
	}

	for _, cartItem := range cartItems {
		cart.Books = append(cart.Books, cartItem)
		cart.TotalPrice += cartItem.Price
	}

	errors := utils.ValidateStruct(cart)

	if errors != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(errors)
	}

	db.Save(cart)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"cart": "cart updated successfully",
	})
}

func RemoveCart(c *fiber.Ctx) error {
	return nil
}
