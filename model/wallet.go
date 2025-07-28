package model

type Wallet struct {
	ID            string `gorm:"type:char(36);primaryKey"`
	Email         string `gorm:"not null;unique"`
	WalletAddress string `gorm:"not null;unique"`
	WalletSecret  string `gorm:"type:text;not null"`
}
