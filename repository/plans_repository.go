package repository

import (
	"github.com/thebytearray/BytePayments/internal/database"
	"github.com/thebytearray/BytePayments/model"
)

func GetPlans() ([]model.Plan, error) {
	var plans []model.Plan

	result := database.DB.Find(&plans)

	return plans, result.Error
}
