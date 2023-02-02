package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nesarptr/book-store-go/config"
)

func GetPK(c fiber.Ctx) error {
	pk, err := config.GetEnv("STRIPE_PUBLISHABLE_KEY")
	if err != nil {
		return fiber.ErrInternalServerError
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"publishableKey": pk,
	})
}
