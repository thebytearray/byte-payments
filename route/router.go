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

	//admin routes
	//
	v1_admin := v1.Group("/admin")
	{
		v1_admin.Post("/login", controller.AdminLoginHandler)
		
		// Protected admin routes
		v1_admin.Use(controller.AdminAuthMiddleware())
		v1_admin.Post("/change-password", controller.ChangePasswordHandler)
		// Plans
		v1_admin.Post("/plans", controller.CreatePlanHandler)
		v1_admin.Put("/plans/:id", controller.UpdatePlanHandler)
		v1_admin.Delete("/plans/:id", controller.DeletePlanHandler)
		// Payments
		v1_admin.Get("/payments", controller.GetAllPaymentsHandler)
		v1_admin.Delete("/payments/:id", controller.DeletePaymentHandler)
		// Wallets
		v1_admin.Get("/wallets", controller.GetAllWalletsHandler)
		v1_admin.Delete("/wallets/:id", controller.DeleteWalletHandler)
		// Currencies
		v1_admin.Post("/currencies", controller.CreateCurrencyHandler)
		v1_admin.Put("/currencies/:code", controller.UpdateCurrencyHandler)
		v1_admin.Delete("/currencies/:code", controller.DeleteCurrencyHandler)
	}

	return app
}
