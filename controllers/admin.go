package controllers

import (
	"fmt"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/nesarptr/book-store-go/config"
	"github.com/nesarptr/book-store-go/models"
	"github.com/nesarptr/book-store-go/utils"
)

func CreateBook(c *fiber.Ctx) error {
	book := new(models.Book)

	book.Title = c.FormValue("name")
	price, err := strconv.ParseFloat(c.FormValue("price"), 32)
	if err != nil {
		return fiber.ErrUnprocessableEntity
	}
	book.Price = float32(price)
	book.Description = c.FormValue("description")
	bookImg, err := c.FormFile("image")
	if err != nil {
		return fiber.ErrUnprocessableEntity
	}

	extension := filepath.Ext(bookImg.Filename)

	validExtensions := []string{".jpg", ".jpeg", ".png"}

	// Check if the extension is valid
	valid := false
	for _, v := range validExtensions {
		if v == extension {
			valid = true
			break
		}
	}

	if !valid {
		return fiber.ErrUnprocessableEntity
	}

	imgUrl := fmt.Sprintf("%d-%s", time.Now().Unix(), bookImg.Filename)
	imgDir := fmt.Sprintf("./images/%s", imgUrl)

	if err := c.SaveFile(bookImg, imgDir); err != nil {
		return fiber.ErrUnprocessableEntity
	}

	book.ImgUrl = imgUrl

	userId := c.Locals("userId").(float64)

	book.UserID = uint(userId)

	errors := utils.ValidateStruct(*book)

	if errors != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(errors)
	}

	db := config.GetDB()

	if err := book.Create(db); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(err.Error())
	}

	user := new(models.User)
	db.First(user, book.UserID)
	user.Books = append(user.Books, *book)
	db.Save(user)
	return c.Status(fiber.StatusCreated).JSON(book)
}

func GetBooks(c *fiber.Ctx) error {
	userId := c.Locals("userId").(float64)
	db := config.GetDB()
	user := new(models.User)
	db.Model(&models.User{}).Order("ID desc").Preload("Books").First(user, userId)
	if len(user.Books) > 0 {
		return c.Status(fiber.StatusOK).JSON(user.Books)
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
	if book.ID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid book id",
		})
	}
	if book.UserID == uint(userId) {
		return c.Status(fiber.StatusOK).JSON(book)
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
	if book.ID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid book id",
		})
	}
	if book.UserID == uint(userId) {
		book.Title = c.FormValue("name")
		price, err := strconv.ParseFloat(c.FormValue("price"), 32)
		if err != nil {
			return fiber.ErrUnprocessableEntity
		}
		book.Price = float32(price)
		book.Description = c.FormValue("description")
		bookImg, err := c.FormFile("image")
		if err == nil {
			extension := filepath.Ext(bookImg.Filename)

			validExtensions := []string{".jpg", ".jpeg", ".png"}

			// Check if the extension is valid
			valid := false
			for _, v := range validExtensions {
				if v == extension {
					valid = true
					break
				}
			}

			if !valid {
				return fiber.ErrUnprocessableEntity
			}

			imgUrl := fmt.Sprintf("%d-%s", time.Now().Unix(), bookImg.Filename)
			imgDir := fmt.Sprintf("./images/%s", imgUrl)

			if err := c.SaveFile(bookImg, imgDir); err != nil {
				return fiber.ErrUnprocessableEntity
			}
			done := make(chan error)

			go utils.RemoveImage(book.ImgUrl, 60, done)
			err := <-done
			if err != nil {
				fmt.Println(err.Error())
			}

			book.ImgUrl = imgUrl

		}
		errors := utils.ValidateStruct(*book)

		if errors != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(errors)
		}
		db.Save(book)

		return c.Status(fiber.StatusOK).JSON(book)
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
	if book.ID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid book id",
		})
	}
	if book.UserID == uint(userId) {
		db.Unscoped().Where("book_id = ?", bookId).Delete(models.CartItem{})
		db.Unscoped().Delete(models.Book{}, bookId)
		done := make(chan error)

		go utils.RemoveImage(book.ImgUrl, 60, done)
		err := <-done
		if err != nil {
			fmt.Println(err.Error())
		}
		return c.Status(fiber.StatusOK).JSON(book)
	}
	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"message": "this book does not belong to user",
	})
}
