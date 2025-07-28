package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/thebytearray/BytePayments/dto"
	"github.com/thebytearray/BytePayments/internal/database"
	"github.com/thebytearray/BytePayments/repository"
	"github.com/thebytearray/BytePayments/service"
)

func GetCurrenciesHandler(ctx *fiber.Ctx) error {
	currencyService := service.NewCurrenciesService(repository.NewCurrenciesRepository(database.DB))
	currencies, err := currencyService.GetCurrencies()
	if err != nil {
		return ctx.JSON(dto.NewError("Currencies not found", err))
	}
	return ctx.JSON(dto.NewSuccess("Currencies fetched successfully", currencies))
}
