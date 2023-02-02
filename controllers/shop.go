package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nesarptr/book-store-go/config"
	"github.com/nesarptr/book-store-go/models"
	"github.com/nesarptr/book-store-go/utils"
)

func GetAllBooks(c *fiber.Ctx) error {
	db := config.GetDB()
	books := new([]models.Book)
	db.Order("ID desc").Find(books)
	if len(*books) > 0 {
		return c.Status(fiber.StatusOK).JSON(books)
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
	return c.Status(fiber.StatusOK).JSON(book)
}

func GetCart(c *fiber.Ctx) error {
	userId := c.Locals("userId").(float64)
	db := config.GetDB()
	cart := new(models.Cart)
	db.Order("ID desc").Preload("Books.Book").Where("user_id = ?", userId).First(cart)
	if cart.ID != 0 {
		return c.Status(fiber.StatusOK).JSON(cart)
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
		ID       uint `json:"id"`
		Quantity int  `json:"quantity"`
	}
	type inputCart struct {
		Books []cartBook `json:"books"`
	}
	InputCart := new(inputCart)
	if err := c.BodyParser(InputCart); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(err.Error())
	}
	cart := new(models.Cart)
	db.Where("user_id = ?", userId).First(cart)
	if cart.ID != 0 {
		db.Unscoped().Where("cart_id = ?", cart.ID).Delete(&models.CartItem{})
		db.Unscoped().Delete(cart)
	}
	if len(InputCart.Books) <= 0 {
		return c.Status(200).JSON(fiber.Map{
			"message": "cart is successfully removed",
		})
	}

	cartItems := make([]models.CartItem, 0)

	cartBooks := map[uint]cartBook{}

	for _, book := range InputCart.Books {
		cartBooks[book.ID] = book
	}

	cart.UserID = uint(userId)
	cart.TotalPrice = 0
	cart.Create(db)

	for _, book := range cartBooks {
		cartItem := new(models.CartItem)
		bookId := book.ID
		cartItem.CartID = cart.ID
		cartItem.BookID = uint(bookId)
		cartItem.Quantity = book.Quantity
		errors := utils.ValidateStruct(cartItem)
		if errors != nil {
			db.Unscoped().Where("cart_id = ?", cart.ID).Delete(&models.CartItem{})
			db.Unscoped().Delete(cart)
			return c.Status(fiber.StatusUnprocessableEntity).JSON(errors)
		}
		cartItems = append(cartItems, *cartItem)
	}

	if db.Create(&cartItems).Error != nil {
		db.Unscoped().Where("cart_id = ?", cart.ID).Delete(&models.CartItem{})
		db.Unscoped().Delete(cart)
		return fiber.ErrBadRequest
	}

	for _, cartItem := range cartItems {
		cart.Books = append(cart.Books, cartItem)
		cart.TotalPrice += cartItem.Price
	}

	errors := utils.ValidateStruct(cart)

	if errors != nil {
		db.Unscoped().Where("cart_id = ?", cart.ID).Delete(&models.CartItem{})
		db.Unscoped().Delete(cart)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(errors)
	}

	db.Save(cart)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"cart": "cart updated successfully",
	})
}

func RemoveCart(c *fiber.Ctx) error {
	userId := c.Locals("userId").(float64)
	cart := new(models.Cart)
	db := config.GetDB()
	db.Where("user_id = ?", userId).First(cart)
	if cart.ID != 0 {
		db.Unscoped().Where("cart_id = ?", cart.ID).Delete(&models.CartItem{})
		db.Unscoped().Delete(cart)
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "cart successfully deleted",
	})
}

func PostOrder(c *fiber.Ctx) error {
	userId := c.Locals("userId").(float64)
	cart := new(models.Cart)
	db := config.GetDB()
	db.Order("ID desc").Preload("Books.Book").Where("user_id = ?", userId).First(cart)
	if cart.ID == 0 {
		return fiber.ErrBadRequest
	}
	order := new(models.Order)
	order.IsPaid = false
	order.UserID = uint(userId)
	order.TotalPrice = cart.TotalPrice
	order.PaymentID = ""
	if db.Create(order).Error != nil {
		return fiber.ErrInternalServerError
	}
	orderItems := make([]models.OrderItem, 0)
	for _, cartItem := range cart.Books {
		orderItem := new(models.OrderItem)
		orderItem.OrderID = order.ID
		orderItem.Title = cartItem.Book.Title
		orderItem.Price = float64(cartItem.Book.Price)
		orderItem.ImgUrl = cartItem.Book.ImgUrl
		orderItem.Description = cartItem.Book.Description
		orderItem.BookID = cartItem.Book.ID
		orderItem.UserID = cartItem.Book.ID
		orderItem.Quantity = cartItem.Quantity
		errors := utils.ValidateStruct(orderItem)
		if errors != nil {
			db.Unscoped().Where("order_id = ?", order.ID).Delete(&models.OrderItem{})
			db.Unscoped().Delete(order)
			return c.Status(fiber.StatusUnprocessableEntity).JSON(errors)
		}
		orderItems = append(orderItems, *orderItem)
	}
	if db.Create(&orderItems).Error != nil {
		db.Unscoped().Where("order_id = ?", order.ID).Delete(&models.OrderItem{})
		db.Unscoped().Delete(order)
		return fiber.ErrBadRequest
	}
	order.Books = orderItems
	errors := utils.ValidateStruct(order)
	if errors != nil {
		db.Unscoped().Where("order_id = ?", order.ID).Delete(&models.OrderItem{})
		db.Unscoped().Delete(order)
		return c.Status(fiber.StatusUnprocessableEntity).JSON(errors)
	}
	if db.Save(order).Error != nil {
		db.Unscoped().Where("order_id = ?", order.ID).Delete(&models.OrderItem{})
		db.Unscoped().Delete(order)
		return fiber.ErrBadRequest
	}
	db.Unscoped().Where("cart_id = ?", cart.ID).Delete(&models.CartItem{})
	db.Unscoped().Delete(cart)
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "order placed successfully"})
}
