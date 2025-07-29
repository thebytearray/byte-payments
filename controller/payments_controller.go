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

// CreatePaymentHandler godoc
// @Summary      Create a payment
// @Description  Creates a payment in the database and returns payment details
// @Tags         create
// @Produce      json
// @Param        request body dto.CreatePaymentRequest true "Payment Request"
// @Success      200  {object}  dto.PaymentResponse
// @Failure      404  {object}  dto.PaymentResponse
// @Router       /api/v1/payments/create [post]
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

	resp, err := paymentService.CreatePayment(body)

	if err != nil {
		return ctx.Status(http.StatusExpectationFailed).JSON(dto.NewError("Payment creation failed", err))
	}

	return ctx.Status(http.StatusOK).JSON(dto.NewSuccess("Payment created", resp))
}

// CancelPaymentHandler godoc
// @Summary      Cancel a created payment
// @Description  Cancels a payment that has been created before
// @Tags         cancel
// @Param        id path string true "Payment ID"
// @Produce      json
// @Success      200  {object}  dto.PaymentResponse
// @Failure      404  {object}  dto.PaymentResponse
// @Router       /api/v1/payments/{id}/cancel [patch]
func CancelPaymentHandler(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	paymentService := service.NewPaymentService(repository.NewPaymentRepository(database.DB))

	resp := paymentService.CancelPaymentById(id)
	return ctx.JSON(resp)
}

// CancelPaymentHandler godoc
// @Summary      Get Status of a payment
// @Description  Gives the status of a payment
// @Tags         status
// @Param        id path string true "Payment ID"
// @Produce      json
// @Success      200  {object}  dto.PaymentResponse
// @Failure      404  {object}  dto.PaymentResponse
// @Router       /api/v1/payments/{id}/status [get]
func GetPaymentStatusHandler(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	paymentService := service.NewPaymentService(repository.NewPaymentRepository(database.DB))

	resp := paymentService.CheckPaymentStatusById(id)
	return ctx.JSON(resp)
}
