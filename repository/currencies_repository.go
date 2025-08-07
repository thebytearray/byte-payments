package repository

import (
	"github.com/thebytearray/BytePayments/model"
	"gorm.io/gorm"
)

type CurrenciesRepository interface {
	GetCurrencies() ([]model.Currency, error)
	GetCurrencyByCode(code string) (*model.Currency, error)
	CreateCurrency(currency *model.Currency) error
	UpdateCurrency(currency *model.Currency) error
	DeleteCurrency(code string) error
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

func (r *currenciesRepository) GetCurrencyByCode(code string) (*model.Currency, error) {
	var currency model.Currency
	res := r.db.First(&currency, "code = ?", code)
	if res.Error != nil {
		return nil, res.Error
	}
	return &currency, nil
}

func (r *currenciesRepository) CreateCurrency(currency *model.Currency) error {
	return r.db.Create(currency).Error
}

func (r *currenciesRepository) UpdateCurrency(currency *model.Currency) error {
	return r.db.Save(currency).Error
}

func (r *currenciesRepository) DeleteCurrency(code string) error {
	return r.db.Delete(&model.Currency{}, "code = ?", code).Error
}
