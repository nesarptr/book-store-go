package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/nesarptr/book-store-go/config"
	"github.com/nesarptr/book-store-go/models"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/paymentintent"
)

func GetPK(c *fiber.Ctx) error {
	pk, err := config.GetEnv("STRIPE_PUBLISHABLE_KEY")
	if err != nil {
		return fiber.ErrInternalServerError
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"key": pk,
	})
}

func Pay(c *fiber.Ctx) error {
	sk, _ := config.GetEnv("STRIPE_KEY")
	stripe.Key = sk
	userId := c.Locals("userId").(float64)
	orderId := c.Params("id")
	order := new(models.Order)
	db := config.GetDB()
	db.First(order, orderId)
	if order.ID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid order id",
		})
	}
	if order.UserID != uint(userId) {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "this order does not belong to user",
		})
	}
	if order.IsPaid {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "the amount is already paid for this order",
		})
	}
	var pi *stripe.PaymentIntent
	if order.PaymentID != "" {
		pi, _ = paymentintent.Get(order.PaymentID, nil)
	} else {
		params := &stripe.PaymentIntentParams{
			Amount:   stripe.Int64(int64(order.TotalPrice * 100)),
			Currency: stripe.String(string(stripe.CurrencyUSD)),
			AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
				Enabled: stripe.Bool(true),
			},
			ReceiptEmail: stripe.String(c.Locals("email").(string)),
		}
		pi, _ = paymentintent.New(params)
		order.PaymentID = pi.ID
		db.Save(order)
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"clientSecret": pi.ClientSecret,
	})
}

func ConfirmPay(c *fiber.Ctx) error {
	sk, _ := config.GetEnv("STRIPE_KEY")
	stripe.Key = sk
	orderId := c.Params("id")
	order := new(models.Order)
	db := config.GetDB()
	db.First(order, orderId)
	if order.ID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid order id",
		})
	}
	pi, _ := paymentintent.Get(order.PaymentID, nil)
	if pi.Status != stripe.PaymentIntentStatusSucceeded {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "user did not pay",
		})
	}

	order.IsPaid = true
	db.Save(order)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "payment successful",
	})

}
