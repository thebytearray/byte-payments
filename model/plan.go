package model

type Plan struct {
	ID           string  `gorm:"type:char(36);primaryKey" json:"id"`
	Name         string  `gorm:"not null" json:"name"`
	Description  string  `gorm:"type:text" json:"description"`
	PriceUSD     float64 `gorm:"not null" json:"price_usd"`
	DurationDays int64   `gorm:"not null" json:"duration_days"`
}
