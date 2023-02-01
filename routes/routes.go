package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nesarptr/book-store-go/controllers"
	"github.com/nesarptr/book-store-go/middleware"
	// "github.com/nesarptr/book-store-go/middleware"
)

func SetUpRoutes(app *fiber.App) {

	// Authentication Routes

	auth := app.Group("/auth")
	auth.Post("/signup", controllers.SignUp)
	auth.Post("/signin", controllers.SignIn)

	// Admin Routes

	admin := app.Group("/admin", middleware.Protected()...)

	// Shop Routes

	// shop := app.Group("/shop", middleware.Protected()...)

}
