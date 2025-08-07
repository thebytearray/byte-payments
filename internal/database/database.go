package database

import (
	"fmt"
	"log"

	"github.com/dgraph-io/ristretto"
	"github.com/thebytearray/BytePayments/config"
	"github.com/thebytearray/BytePayments/internal/util"
	"github.com/thebytearray/BytePayments/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB
var Cache *ristretto.Cache

func Connect() {
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Cfg.DATABASE_USER,
		config.Cfg.DATABASE_PASS,
		config.Cfg.DATABASE_HOST,
		config.Cfg.DATABASE_PORT,
		config.Cfg.DATABASE_NAME,
	)

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatalln("Could not connect to the database:", err)
	}

	log.Println("Connected to database successfully 📦")

	// Initialize Ristretto cache
	Cache, err = ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
	})
	if err != nil {
		log.Fatalln("Could not initialize Ristretto cache:", err)
	}
	log.Println("Initialized Ristretto cache successfully ⚡")

	//migrate
	//
	//
	//
	err = DB.AutoMigrate(&model.Currency{}, &model.Payment{}, &model.Plan{}, &model.Wallet{}, &model.Admin{})
	if err != nil {
		log.Printf("Failed to automigrate database, %v", err)
	}
}

func SeedDatabase() {

	// Seed default admin
	SeedAdmin()

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

func SeedAdmin() {
	// Check if admin already exists
	var count int64
	DB.Model(&model.Admin{}).Count(&count)
	if count > 0 {
		log.Println("Admin already exists, skipping seeding")
		return
	}

	// Hash default password
	hashedPassword, err := util.HashPassword("admin123")
	if err != nil {
		log.Printf("Failed to hash admin password: %v", err)
		return
	}

	admin := &model.Admin{
		Username: "admin",
		Email:    "admin@bytepayments.com",
		Password: hashedPassword,
		IsActive: true,
	}

	result := DB.Create(admin)
	if result.Error != nil {
		log.Printf("Failed to create admin: %v", result.Error)
	} else {
		log.Println("Default admin created successfully (username: admin, password: admin123)")
	}
}
