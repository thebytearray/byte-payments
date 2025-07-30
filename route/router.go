package route

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/thebytearray/BytePayments/config"
	"github.com/thebytearray/BytePayments/controller"
)

func NewRouter() *fiber.App {
	app := fiber.New()

	// Enable CORS for frontend integration
	app.Use(func(c *fiber.Ctx) error {
		c.Set("Access-Control-Allow-Origin", "*")
		c.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
		c.Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")
		
		if c.Method() == "OPTIONS" {
			return c.SendStatus(200)
		}
		
		return c.Next()
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
