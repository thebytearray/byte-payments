package model

import "time"

type Wallet struct {
	ID            string    `gorm:"type:char(27);primaryKey" json:"id"`
	Email         string    `gorm:"not null;unique" json:"email"`
	WalletAddress string    `gorm:"not null;unique" json:"tron_address"`
	WalletSecret  string    `gorm:"type:text;not null" json:"-"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
