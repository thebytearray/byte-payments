package model

type Plan struct {
	ID           string  `gorm:"type:char(36);primaryKey"`
	Name         string  `gorm:"not null"`
	Description  string  `gorm:"type:text"`
	PriceUSD     float64 `gorm:"not null"`
	DurationDays int64   `gorm:"not null"`
}
