package model

type Currency struct {
	Code         string `gorm:"primaryKey;size:10"`
	Name         string `gorm:"not null;size:50"`
	Network      string `gorm:"not null;size:20"`
	IsToken      bool   `gorm:"default:false"`
	ContractAddr string `gorm:"size:50"`
	Enabled      bool   `gorm:"default:true"`
}
