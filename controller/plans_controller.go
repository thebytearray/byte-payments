package controller

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/thebytearray/BytePayments/dto"
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

	plans, err := service.GetPlans()

	if err != nil {
		return ctx.JSON(
			dto.ApiResponse{
				Status:     "error",
				StatusCode: http.StatusNotFound,
				Message:    "Plan not found",
				Data:       nil,
				Error:      err.Error(),
			},
		)
	}

	return ctx.JSON(
		dto.ApiResponse{
			Status:     "ok",
			StatusCode: http.StatusOK,
			Message:    "Plans fetched successfully",
			Data:       plans,
		},
	)
}
