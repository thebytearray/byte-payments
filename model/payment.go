package model

import "time"

type PaymentStatus string

const (
	Pending   PaymentStatus = "pending"
	Completed PaymentStatus = "completed"
	Cancelled PaymentStatus = "cancelled"
)

type Payment struct {
	ID     string `gorm:"type:char(36);primaryKey"`
	PlanID string `gorm:"not null"`          // FK field
	Plan   Plan   `gorm:"foreignKey:PlanID"` // Assoc

	WalletID string `gorm:"not null"`            // FK field
	Wallet   Wallet `gorm:"foreignKey:WalletID"` // Assoc

	CurrencyCode string   `gorm:"size:10;not null"`                        // FK field
	Currency     Currency `gorm:"foreignKey:CurrencyCode;references:Code"` // Assoc

	AmountUSD float64 `gorm:"not null"`
	AmountTRX float64 `gorm:"not null"`
	UserEmail string  `gorm:"not null"`

	Status        PaymentStatus `gorm:"type:varchar(20);default:'pending'"` // enum-like string
	PaidAmountTRX float64       `gorm:"default:0"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
