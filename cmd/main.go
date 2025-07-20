package main

import (
	"github.com/TheByteArray/BytePayments/config"
	"github.com/TheByteArray/BytePayments/internal/database"
)

func main() {
	config.LoadConfig()
	database.Connect()

}
