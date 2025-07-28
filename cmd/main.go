package main

import (
	"github.com/thebytearray/BytePayments/config"
	_ "github.com/thebytearray/BytePayments/docs"
	"github.com/thebytearray/BytePayments/internal/database"
	"github.com/thebytearray/BytePayments/internal/tron"
	"github.com/thebytearray/BytePayments/route"
)

func main() {
	config.NewConfig()
	database.NewConnection()
	tron.NewClient()
	//	database.SeedDatabase()
	//	log.Println(tron.ConvertUSDToTRX(10.00))
	app := route.NewRouter()
	app.Listen(":8080")
	//
	//

}
