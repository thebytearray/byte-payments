package model

type Currency struct {
	Code         string `gorm:"primaryKey;size:27" json:"code"`
	Name         string `gorm:"not null;size:50" json:"name"`
	Network      string `gorm:"not null;size:20" json:"network"`
	IsToken      bool   `gorm:"default:false" json:"isToken"`
	ContractAddr string `gorm:"size:50" json:"contract_addr"`
	Enabled      bool   `gorm:"default:true" json:"enabled"`
}
