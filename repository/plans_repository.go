package repository

import (
	"github.com/TheByteArray/BytePayments/internal/database"
	"github.com/TheByteArray/BytePayments/model"
)

func GetPlans() ([]model.Plan, error) {
	var plans []model.Plan

	result := database.DB.Find(&plans)

	return plans, result.Error
}
