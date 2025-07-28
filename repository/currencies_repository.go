package repository

import (
	"github.com/thebytearray/BytePayments/model"
	"gorm.io/gorm"
)

type CurrenciesRepository interface {
	GetCurrencies() ([]model.Currency, error)
}

type currenciesRepository struct {
	db *gorm.DB
}

func NewCurrenciesRepository(db *gorm.DB) CurrenciesRepository {
	return &currenciesRepository{db}
}

func (r *currenciesRepository) GetCurrencies() ([]model.Currency, error) {
	var currencies []model.Currency
	res := r.db.Find(&currencies)
	return currencies, res.Error
}
