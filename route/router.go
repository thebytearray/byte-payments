package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/thebytearray/BytePayments/config"
	"github.com/thebytearray/BytePayments/controller"
)

func NewRouter() *fiber.App {
	app := fiber.New()

	// Frontend routes
	app.Get("/", func(c *fiber.Ctx) error {
		return c.Redirect("/plans")
	})

	app.Get("/plans", func(c *fiber.Ctx) error {
		return c.SendFile("./static/select-plan.html")
	})

	app.Get("/select-plan", func(c *fiber.Ctx) error {
		return c.Redirect("/plans")
	})

	app.Get("/pay", func(c *fiber.Ctx) error {
		return c.SendFile("./static/pay.html")
	})

	app.Get("/payment.html", func(c *fiber.Ctx) error {
		return c.Redirect("/pay")
	})

	// API routes - Swagger UI only in development mode
	if config.Cfg.APP_ENV == "development" {
		app.Get("/swagger/*", swagger.HandlerDefault)
	}
	
	v1 := app.Group("/api/v1")
	
	//verification
	v1_verification := v1.Group("/verification")
	{
		v1_verification.Post("/send-code", controller.SendVerificationCodeHandler)
		v1_verification.Post("/verify-code", controller.VerifyEmailCodeHandler)
	}
	
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
