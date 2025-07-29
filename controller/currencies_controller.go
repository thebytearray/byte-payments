package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/thebytearray/BytePayments/dto"
	"github.com/thebytearray/BytePayments/internal/database"
	"github.com/thebytearray/BytePayments/repository"
	"github.com/thebytearray/BytePayments/service"
)

// GetCurrenciesHandler godoc
// @Summary      Get all available currencies
// @Description  Returns a list of all supported currencies with their exchange rates and symbols
// @Tags         currencies
// @Accept       json
// @Produce      json
// @Success      200  {object}  dto.ApiResponse "Currencies retrieved successfully"
// @Failure      404  {object}  dto.ApiResponse "No currencies found"
// @Router       /api/v1/currencies [get]
func GetCurrenciesHandler(ctx *fiber.Ctx) error {
	currencyService := service.NewCurrenciesService(repository.NewCurrenciesRepository(database.DB))
	currencies, err := currencyService.GetCurrencies()
	if err != nil {
		return ctx.JSON(dto.NewError("Currencies not found", err))
	}
	return ctx.JSON(dto.NewSuccess("Currencies fetched successfully", currencies))
}
