package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/helmet/v2"
	"github.com/nesarptr/book-store-go/config"
	"github.com/nesarptr/book-store-go/models"
	"github.com/nesarptr/book-store-go/routes"
)

func main() {
	app := fiber.New()

	app.Use(cors.New())
	app.Use(helmet.New())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	routes.SetUpRoutes(app)

	log.Fatal(app.Listen(":3000"))
}

func init() {
	err := config.Connect()
	if err != nil {
		fmt.Println(err)
	}
	db := config.GetDB()
	db.AutoMigrate(&models.User{}, &models.Book{}, &models.Order{}, &models.Cart{}, &models.CartItem{}, &models.OrderItem{})
	fmt.Println("Database successfully connected!")
}
