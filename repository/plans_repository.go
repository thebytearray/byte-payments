package repository

import (
	"github.com/thebytearray/BytePayments/model"
	"gorm.io/gorm"
)

type PlanRepository interface {
	GetPlans() ([]model.Plan, error)
	GetPlanByID(id string) (*model.Plan, error)
	CreatePlan(plan *model.Plan) error
	UpdatePlan(plan *model.Plan) error
	DeletePlan(id string) error
}

type planRepository struct {
	db *gorm.DB
}

func NewPlansRepository(db *gorm.DB) PlanRepository {
	return &planRepository{db}
}
func (r *planRepository) GetPlans() ([]model.Plan, error) {
	var plans []model.Plan
	res := r.db.Find(&plans)
	return plans, res.Error
}

func (r *planRepository) GetPlanByID(id string) (*model.Plan, error) {
	var plan model.Plan
	res := r.db.First(&plan, "id = ?", id)
	if res.Error != nil {
		return nil, res.Error
	}
	return &plan, nil
}

func (r *planRepository) CreatePlan(plan *model.Plan) error {
	return r.db.Create(plan).Error
}

func (r *planRepository) UpdatePlan(plan *model.Plan) error {
	return r.db.Save(plan).Error
}

func (r *planRepository) DeletePlan(id string) error {
	return r.db.Delete(&model.Plan{}, "id = ?", id).Error
}
