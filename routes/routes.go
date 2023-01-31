package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nesarptr/book-store-go/controllers"
)

func SetUpRoutes(app *fiber.App) {

	// Authentication Routes

	auth := app.Group("/auth")
	auth.Post("/signup", controllers.SignUp)

	// Admin Routes
	// Shop Routes
}
