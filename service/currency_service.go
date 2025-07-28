package service

import (
	"github.com/thebytearray/BytePayments/model"
	"github.com/thebytearray/BytePayments/repository"
)

type CurrenciesService interface {
	GetCurrencies() ([]model.Currency, error)
}

type currenciesService struct {
	repo repository.CurrenciesRepository
}

func NewCurrenciesService(repo repository.CurrenciesRepository) CurrenciesService {
	return &currenciesService{repo}
}

func (s *currenciesService) GetCurrencies() ([]model.Currency, error) {
	return s.repo.GetCurrencies()
}
