package controller

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/thebytearray/BytePayments/dto"
	"github.com/thebytearray/BytePayments/service"
)

func GetCurrenciesHandler(ctx *fiber.Ctx) error {
	currencies, err := service.GetCurrencies()

	if err != nil {
		return ctx.JSON(
			dto.ApiResponse{
				Status:     string(dto.ERROR),
				StatusCode: http.StatusNotFound,
				Message:    "Currencies not found",
				Data:       nil,
				Error:      err.Error(),
			},
		)
	}

	return ctx.JSON(
		dto.ApiResponse{
			Status:     string(dto.OK),
			StatusCode: http.StatusOK,
			Message:    "Currencies fetched successfully",
			Data:       currencies,
		},
	)
}
