// Package main BytePayments API
// @title BytePayments API
// @version 1.0
// @description A payment processing API for BytePayments platform
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.email support@bytepayments.com
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @host localhost:8080
// @BasePath /
// @schemes http https
package main

import (
	"github.com/thebytearray/BytePayments/config"
	_ "github.com/thebytearray/BytePayments/docs"
	"github.com/thebytearray/BytePayments/internal/cron"
	"github.com/thebytearray/BytePayments/internal/database"
	"github.com/thebytearray/BytePayments/internal/tron"
	"github.com/thebytearray/BytePayments/route"
)

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	config.NewConfig()
	database.Connect()
	tron.NewClient()
	//database.SeedDatabase()
	//	log.Println(tron.ConvertUSDToTRX(10.00))
	go cron.NewPaymentCron()
	
	// Seed admin before starting server
	database.SeedAdmin()
	
	app := route.NewRouter()
	app.Listen(":8080")

}
