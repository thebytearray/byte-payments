package route

import (
	"github.com/TheByteArray/BytePayments/controller"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
)

func NewRouter() *fiber.App {
	app := fiber.New()
	app.Get("/swagger/*", swagger.HandlerDefault)
	v1 := app.Group("/api/v1")
	//payments
	//
	v1_payments := v1.Group("/payments")
	{
		v1_payments.Post("/create", controller.CreatePaymentHandler)
		v1_payments.Patch("/:id/cancel", controller.CancelPaymentHandler)
		v1_payments.Get("/:id/status", controller.GetPaymentStatusHandler)
	}
	//plans
	//
	v1.Get("/plans", controller.GetPlansHandler)
	v1.Get("/currencies", controller.GetCurrenciesHandler)

	return app
}
