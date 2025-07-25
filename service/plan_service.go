package service

import (
	"github.com/TheByteArray/BytePayments/model"
	"github.com/TheByteArray/BytePayments/repository"
)

func GetPlans() ([]model.Plan, error) {
	return repository.GetPlans()
}
