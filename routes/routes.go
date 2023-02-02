package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nesarptr/book-store-go/controllers"
	"github.com/nesarptr/book-store-go/middleware"
)

func SetUpRoutes(app *fiber.App) {

	// Authentication Routes

	auth := app.Group("/auth")
	auth.Post("/signup", controllers.SignUp)
	auth.Post("/signin", controllers.SignIn)
	auth.Get("/jwt", middleware.Protected()...)
	auth.Get("/jwt", controllers.Jwt)

	// Admin Routes

	admin := app.Group("/admin", middleware.Protected()...)

	admin.Post("/book", controllers.CreateBook)
	admin.Get("/books", controllers.GetBooks)
	admin.Get("/book/:id", controllers.GetBook)
	admin.Put("/book/:id", controllers.UpdateBook)
	admin.Delete("/book/:id", controllers.DeleteBook)

	// Shop Routes

	shop := app.Group("/shop", middleware.Protected()...)

	shop.Get("/books", controllers.GetAllBooks)
	shop.Get("/book/:id", controllers.GetSingleBook)
	shop.Get("/cart", controllers.GetCart)
	shop.Put("/cart", controllers.PostCart)
	shop.Delete("/cart", controllers.RemoveCart)
	shop.Post("/order", controllers.PostOrder)
}
