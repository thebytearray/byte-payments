package database

import (
	"fmt"
	"log"

	"github.com/TheByteArray/BytePayments/config"
	"github.com/TheByteArray/BytePayments/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		config.Cfg.DATABASE_USER,
		config.Cfg.DATABASE_PASS,
		config.Cfg.DATABASE_HOST,
		config.Cfg.DATABASE_PORT,
		config.Cfg.DATABASE_NAME)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Printf("Failed to connect to database, %v", err)
	}

	//migrate
	//
	//
	//
	err = db.AutoMigrate(&model.Currency{}, &model.Payment{}, &model.Plan{})
	if err != nil {
		log.Printf("Failed to automigrate database, %v", err)
	}

	DB = db

	log.Println("Connected to database successfully.")

}
