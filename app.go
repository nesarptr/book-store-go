package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/helmet/v2"
	"github.com/nesarptr/book-store-go/config"
	"github.com/nesarptr/book-store-go/models"
	"github.com/nesarptr/book-store-go/routes"
)

func main() {
	if err := os.MkdirAll(filepath.Dir("./images/"), os.ModePerm); err != nil {
		fmt.Println(err.Error())
	}

	port := config.GetEnv("PORT")

	if port == "" {
		port = "4000"
	}

	app := fiber.New()
	app.Use(cors.New())
	app.Use(helmet.New())
	app.Use(compress.New())
	app.Static("/", "./images")

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	routes.SetUpRoutes(app)

	fmt.Println(app.Listen(":" + port))
}

func init() {
	err := config.Connect()
	if err != nil {
		fmt.Println(err.Error())
	}
	db := config.GetDB()
	db.AutoMigrate(&models.User{}, &models.Book{}, &models.Order{}, &models.Cart{}, &models.CartItem{}, &models.OrderItem{})
	fmt.Println("Database successfully connected!")
}
