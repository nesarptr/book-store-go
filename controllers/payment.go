package controllers

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/nesarptr/book-store-go/config"
	"github.com/nesarptr/book-store-go/models"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/paymentintent"
)

func GetPK(c *fiber.Ctx) error {
	pk := config.GetEnv("STRIPE_PUBLISHABLE_KEY")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"key": pk,
	})
}

func Pay(c *fiber.Ctx) error {
	sk := config.GetEnv("STRIPE_KEY")
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
	var err error
	if order.PaymentID != "" {
		pi, err = paymentintent.Get(order.PaymentID, nil)
		if err != nil {
			log.Fatal(err.Error())
			return fiber.ErrInternalServerError
		}
	} else {
		params := &stripe.PaymentIntentParams{
			Amount:   stripe.Int64(int64(order.TotalPrice * 100)),
			Currency: stripe.String(string(stripe.CurrencyUSD)),
			AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
				Enabled: stripe.Bool(true),
			},
			ReceiptEmail: stripe.String(c.Locals("email").(string)),
		}
		pi, err = paymentintent.New(params)
		if err != nil {
			log.Fatal(err.Error())
			return fiber.ErrInternalServerError
		}
		order.PaymentID = pi.ID
		db.Save(order)
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"clientSecret": pi.ClientSecret,
	})
}

func ConfirmPay(c *fiber.Ctx) error {
	sk := config.GetEnv("STRIPE_KEY")
	stripe.Key = sk
	paymentId := c.Params("id")
	order := new(models.Order)
	db := config.GetDB()
	db.Where("payment_id = ?", paymentId).First(order)
	if order.ID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "invalid order",
		})
	}
	pi, err := paymentintent.Get(paymentId, nil)
	if err != nil {
		log.Fatal(err.Error())
		return fiber.ErrInternalServerError
	}
	if pi.Status != stripe.PaymentIntentStatusSucceeded {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "user did not pay",
		})
	}

	order.IsPaid = true
	db.Save(order)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "payment successful",
		"id": order.ID,
	})

}
