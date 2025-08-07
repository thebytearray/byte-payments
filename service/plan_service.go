package service

import (
	"github.com/thebytearray/BytePayments/model"
	"github.com/thebytearray/BytePayments/repository"
)

type PlanService interface {
	GetPlans() ([]model.Plan, error)
	GetPlanByID(id string) (*model.Plan, error)
	CreatePlan(plan *model.Plan) error
	UpdatePlan(plan *model.Plan) error
	DeletePlan(id string) error
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

func (s *plansService) GetPlanByID(id string) (*model.Plan, error) {
	return s.repo.GetPlanByID(id)
}

func (s *plansService) CreatePlan(plan *model.Plan) error {
	return s.repo.CreatePlan(plan)
}

func (s *plansService) UpdatePlan(plan *model.Plan) error {
	return s.repo.UpdatePlan(plan)
}

func (s *plansService) DeletePlan(id string) error {
	return s.repo.DeletePlan(id)
}
