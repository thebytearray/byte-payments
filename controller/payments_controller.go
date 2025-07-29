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
// @Summary      Create a new payment
// @Description  Creates a new payment with the specified plan and currency, generates a TRX wallet address and QR code for payment
// @Tags         payments
// @Accept       json
// @Produce      json
// @Param        request body dto.CreatePaymentRequest true "Payment creation request"
// @Success      200  {object}  dto.ApiResponse{data=dto.PaymentResponse} "Payment created successfully"
// @Failure      400  {object}  dto.ApiResponse "Invalid request body or validation error"
// @Failure      422  {object}  dto.ApiResponse "Payment creation failed"
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
// @Summary      Cancel a payment
// @Description  Cancels an existing payment by its ID
// @Tags         payments
// @Accept       json
// @Produce      json
// @Param        id path string true "Payment ID" example("payment_123456")
// @Success      200  {object}  dto.ApiResponse "Payment cancelled successfully"
// @Failure      404  {object}  dto.ApiResponse "Payment not found"
// @Router       /api/v1/payments/{id}/cancel [patch]
func CancelPaymentHandler(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	paymentService := service.NewPaymentService(repository.NewPaymentRepository(database.DB))

	resp := paymentService.CancelPaymentById(id)
	return ctx.JSON(resp)
}

// GetPaymentStatusHandler godoc
// @Summary      Get payment status
// @Description  Retrieves the current status of a payment by its ID
// @Tags         payments
// @Accept       json
// @Produce      json
// @Param        id path string true "Payment ID" example("payment_123456")
// @Success      200  {object}  dto.ApiResponse{data=dto.PaymentResponse} "Payment status retrieved successfully"
// @Failure      404  {object}  dto.ApiResponse "Payment not found"
// @Router       /api/v1/payments/{id}/status [get]
func GetPaymentStatusHandler(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	paymentService := service.NewPaymentService(repository.NewPaymentRepository(database.DB))

	resp := paymentService.CheckPaymentStatusById(id)
	return ctx.JSON(resp)
}
