package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	//app related stuff
	APP_NAME string
	APP_ENV  string
	APP_PORT string
	APP_URL  string
	// database related stuff
	DATABASE_NAME string
	DATABASE_HOST string
	DATABASE_USER string
	DATABASE_PORT string
	DATABASE_PASS string
	//  wallet stuff
	TRX_HOT_WALLET_ADDRESS    string
	TRX_WALLET_ENCRYPTION_KEY string
	TRON_GRID_API_KEY         string
	TRON_GRPC_MAINNET         string
	TRON_GRPC_TESTNET         string
	//emailing config stuff
	EMAIL_SMTP_HOST string
	EMAIL_SMTP_PORT int
	EMAIL_USERNAME  string
	EMAIL_PASSWORD  string
	EMAIL_FROM_NAME string
	EMAIL_FROM_ADDR string
}

var Cfg *Config

func NewConfig() {
	// load the env
	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("Failed to load env, %v", err)
	}

	port, err := strconv.Atoi(os.Getenv("EMAIL_SMTP_PORT"))

	if err != nil {
		port = 587
	}

	Cfg = &Config{
		APP_NAME: os.Getenv("APP_NAME"),
		APP_ENV:  os.Getenv("APP_ENV"),
		APP_PORT: os.Getenv("APP_PORT"),
		APP_URL:  os.Getenv("APP_URL"),

		DATABASE_NAME: os.Getenv("DATABASE_NAME"),
		DATABASE_HOST: os.Getenv("DATABASE_HOST"),
		DATABASE_USER: os.Getenv("DATABASE_USER"),
		DATABASE_PORT: os.Getenv("DATABASE_PORT"),
		DATABASE_PASS: os.Getenv("DATABASE_PASS"),

		TRX_HOT_WALLET_ADDRESS:    os.Getenv("TRX_HOT_WALLET_ADDRESS"),
		TRX_WALLET_ENCRYPTION_KEY: os.Getenv("TRX_WALLET_ENCRYPTION_KEY"),
		TRON_GRID_API_KEY:         os.Getenv("TRON_GRID_API_KEY"),
		TRON_GRPC_MAINNET:         os.Getenv("TRON_GRPC_MAINNET"),
		TRON_GRPC_TESTNET:         os.Getenv("TRON_GRPC_TESTNET"),

		EMAIL_SMTP_HOST: os.Getenv("EMAIL_SMTP_HOST"),
		EMAIL_SMTP_PORT: port,
		EMAIL_USERNAME:  os.Getenv("EMAIL_USERNAME"),
		EMAIL_PASSWORD:  os.Getenv("EMAIL_PASSWORD"),
		EMAIL_FROM_NAME: os.Getenv("EMAIL_FROM_NAME"),
		EMAIL_FROM_ADDR: os.Getenv("EMAIL_FROM_ADDR"),
	}

}
