package repository

import (
	"github.com/thebytearray/BytePayments/model"
	"gorm.io/gorm"
)

type PlanRepository interface {
	GetPlans() ([]model.Plan, error)
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
