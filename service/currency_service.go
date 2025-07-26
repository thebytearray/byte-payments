package service

import (
	"github.com/thebytearray/BytePayments/model"
	"github.com/thebytearray/BytePayments/repository"
)

func GetCurrencies() ([]model.Currency, error) {

	return repository.GetCurrencies()
}
