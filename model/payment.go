package model

import "time"

type PaymentStatus string

const (
	Pending   PaymentStatus = "pending"
	Completed PaymentStatus = "completed"
	Cancelled PaymentStatus = "cancelled"
)

type Payment struct {
	ID            string        `gorm:"type:char(36);primaryKey"`
	PlanID        string        `gorm:"not null"`
	Plan          Plan          `gorm:"foreignKey:PlanID"`
	CurrencyCode  string        `gorm:"size:10;not null"`
	Currency      Currency      `gorm:"foreignKey:CurrencyCode"`
	AmountUSD     float64       `gorm:"not null"`
	AmountTRX     float64       `gorm:"not null"`
	WalletAddress string        `gorm:"not null;unique"`
	WalletSecret  string        `gorm:"type:text;not null"`
	UserEmail     string        `gorm:"not null"`
	Status        PaymentStatus `gorm:"type:enum('pending','completed','cancelled');default:'pending'"`
	PaidAmountTRX float64       `gorm:"default:0"`
	IsSwept       bool          `gorm:"default:false"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
