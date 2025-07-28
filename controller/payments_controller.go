package controller

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/thebytearray/BytePayments/dto"
	"github.com/thebytearray/BytePayments/internal/database"
	"github.com/thebytearray/BytePayments/repository"
	"github.com/thebytearray/BytePayments/service"
)

func CreatePaymentHandler(ctx *fiber.Ctx) error {
	var body dto.CreatePaymentRequest
	//validate body struct
	if err := ctx.BodyParser(&body); err != nil {
		return ctx.JSON(dto.NewError(
			"Invalid request body :(", err))
	}
	//validate all the fields
	validate := validator.New()
	if err := validate.Struct(body); err != nil {

		return ctx.JSON(
			dto.NewError("All fields are required :(", err))

	}

	//payment service
	//
	paymentService := service.NewPaymentService(repository.NewPaymentRepository(database.DB))

	resp, err := paymentService.CreatePayment(ctx.Context(), body)

	if err != nil {
		return ctx.Status(http.StatusExpectationFailed).JSON(dto.NewError("Payment creation failed", err))
	}

	return ctx.Status(http.StatusOK).JSON(dto.NewSuccess("Payment created", resp))
}

func CancelPaymentHandler(ctx *fiber.Ctx) error {

	return nil
}

func GetPaymentStatusHandler(ctx *fiber.Ctx) error {

	return nil
}
