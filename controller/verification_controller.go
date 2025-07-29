package controller

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/thebytearray/BytePayments/dto"
	"github.com/thebytearray/BytePayments/service"
)

// SendVerificationCodeHandler godoc
// @Summary      Send email verification code
// @Description  Sends a verification code to the provided email address for account verification
// @Tags         verification
// @Accept       json
// @Produce      json
// @Param        request body dto.SendVerificationCodeRequest true "Email verification request"
// @Success      200  {object}  dto.ApiResponse{data=dto.SendVerificationCodeResponse} "Verification code sent successfully"
// @Failure      400  {object}  dto.ApiResponse "Invalid email format or request body"
// @Failure      500  {object}  dto.ApiResponse "Failed to send verification code"
// @Router       /api/v1/verification/send-code [post]
func SendVerificationCodeHandler(ctx *fiber.Ctx) error {
	var body dto.SendVerificationCodeRequest
	
	// Parse request body
	if err := ctx.BodyParser(&body); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(dto.NewError(
			"Invalid request body", err))
	}
	
	// Validate request
	validate := validator.New()
	if err := validate.Struct(body); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(
			dto.NewError("Invalid email format", err))
	}
	
	// Send verification code
	verificationService := service.NewVerificationService()
	err := verificationService.GenerateAndSendCode(body.Email)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(
			dto.NewError("Failed to send verification code", err))
	}
	
	return ctx.Status(http.StatusOK).JSON(
		dto.NewSuccess("Verification code sent successfully", 
			dto.SendVerificationCodeResponse{
				Message: "Verification code has been sent to your email address",
			}))
}

// VerifyEmailCodeHandler godoc
// @Summary      Verify email verification code
// @Description  Verifies the email verification code and returns a verification token upon successful verification
// @Tags         verification
// @Accept       json
// @Produce      json
// @Param        request body dto.VerifyEmailCodeRequest true "Email verification code request"
// @Success      200  {object}  dto.ApiResponse{data=dto.VerifyEmailCodeResponse} "Email verified successfully"
// @Failure      400  {object}  dto.ApiResponse "Invalid request body or verification failed"
// @Failure      401  {object}  dto.ApiResponse "Invalid verification code"
// @Router       /api/v1/verification/verify-code [post]
func VerifyEmailCodeHandler(ctx *fiber.Ctx) error {
	var body dto.VerifyEmailCodeRequest
	
	// Parse request body
	if err := ctx.BodyParser(&body); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(dto.NewError(
			"Invalid request body", err))
	}
	
	// Validate request
	validate := validator.New()
	if err := validate.Struct(body); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(
			dto.NewError("Email and code are required", err))
	}
	
	// Verify code
	verificationService := service.NewVerificationService()
	isValid, verificationToken, err := verificationService.VerifyCode(body.Email, body.Code)
	
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(
			dto.NewError("Verification failed", err))
	}
	
	if !isValid {
		return ctx.Status(http.StatusUnauthorized).JSON(
			dto.NewError("Invalid verification code", nil))
	}
	
	return ctx.Status(http.StatusOK).JSON(
		dto.NewSuccess("Email verified successfully", 
			dto.VerifyEmailCodeResponse{
				Message:           "Email verified successfully",
				IsValid:           true,
				VerificationToken: verificationToken,
			}))
} 