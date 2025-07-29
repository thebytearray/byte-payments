package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/thebytearray/BytePayments/dto"
	"github.com/thebytearray/BytePayments/internal/database"
	"github.com/thebytearray/BytePayments/repository"
	"github.com/thebytearray/BytePayments/service"
)

// GetPlansHandler godoc
// @Summary      Get all plans
// @Description  Returns a list of all subscription plans
// @Tags         plans
// @Produce      json
// @Success      200  {object}  dto.ApiResponse
// @Failure      404  {object}  dto.ApiResponse
// @Router       /api/v1/plans [get]
func GetPlansHandler(ctx *fiber.Ctx) error {

	plansService := service.NewPlansService(repository.NewPlansRepository(database.DB))

	plans, err := plansService.GetPlans()

	if err != nil {
		return ctx.JSON(dto.NewError("Plan not found", err))
	}

	return ctx.JSON(dto.NewSuccess("Plans fetched successfully.", plans))
}
