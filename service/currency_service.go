package service

import (
	"github.com/thebytearray/BytePayments/model"
	"github.com/thebytearray/BytePayments/repository"
)

type CurrenciesService interface {
	GetCurrencies() ([]model.Currency, error)
	GetCurrencyByCode(code string) (*model.Currency, error)
	CreateCurrency(currency *model.Currency) error
	UpdateCurrency(currency *model.Currency) error
	DeleteCurrency(code string) error
}

type currenciesService struct {
	repo repository.CurrenciesRepository
}

func NewCurrenciesService(repo repository.CurrenciesRepository) CurrenciesService {
	return &currenciesService{repo}
}

func (s *currenciesService) GetCurrencies() ([]model.Currency, error) {
	currencies, err := s.repo.GetCurrencies()
	if err != nil {
		return nil, err
	}
	
	// Add compatibility fields
	for i := range currencies {
		currencies[i].Symbol = currencies[i].Code
		currencies[i].IsActive = currencies[i].Enabled
	}
	
	return currencies, nil
}

func (s *currenciesService) GetCurrencyByCode(code string) (*model.Currency, error) {
	return s.repo.GetCurrencyByCode(code)
}

func (s *currenciesService) CreateCurrency(currency *model.Currency) error {
	return s.repo.CreateCurrency(currency)
}

func (s *currenciesService) UpdateCurrency(currency *model.Currency) error {
	return s.repo.UpdateCurrency(currency)
}

func (s *currenciesService) DeleteCurrency(code string) error {
	return s.repo.DeleteCurrency(code)
}
