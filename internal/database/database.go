package database

import (
	"fmt"
	"log"

	"github.com/thebytearray/BytePayments/config"
	"github.com/thebytearray/BytePayments/internal/util"
	"github.com/thebytearray/BytePayments/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func NewConnection() {
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

func SeedDatabase() {

	plans := []*model.Plan{
		{
			ID:           util.GenerateUniqueID(),
			Name:         "Basic",
			Description:  "Basic plan",
			PriceUSD:     1.00,
			DurationDays: 7,
		},
		{
			ID:           util.GenerateUniqueID(),
			Name:         "Pro",
			Description:  "Pro plan",
			PriceUSD:     10.00,
			DurationDays: 30,
		},
		{
			ID:           util.GenerateUniqueID(),
			Name:         "Ultimate",
			Description:  "Ultimate plan",
			PriceUSD:     30.00,
			DurationDays: 30,
		},
	}

	currencies := []model.Currency{
		{
			Code:         "TRX",
			Name:         "Tron",
			Network:      "TRC20", // show TRC20 for clarity
			IsToken:      false,
			ContractAddr: "",
			Enabled:      true,
		},
	}
	plansRes := DB.Create(plans)
	currenciesRes := DB.Create(&currencies)
	if plansRes.Error != nil {
		log.Println(plansRes.Error)
	}

	if currenciesRes.Error != nil {
		log.Println(currenciesRes.Error)
	}
}
