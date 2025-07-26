package repository

import (
	"github.com/thebytearray/BytePayments/internal/database"
	"github.com/thebytearray/BytePayments/model"
)

func GetCurrencies() ([]model.Currency, error) {
	var currencies []model.Currency

	result := database.DB.Find(&currencies)

	return currencies, result.Error
}
