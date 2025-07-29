package service

import (
	"github.com/thebytearray/BytePayments/model"
	"github.com/thebytearray/BytePayments/repository"
)

type PlanService interface {
	GetPlans() ([]model.Plan, error)
}

type plansService struct {
	repo repository.PlanRepository
}

func NewPlansService(repo repository.PlanRepository) PlanService {
	return &plansService{repo}
}

func (s *plansService) GetPlans() ([]model.Plan, error) {
	return s.repo.GetPlans()
}
