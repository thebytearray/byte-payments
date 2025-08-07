package controller

import (
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/thebytearray/BytePayments/dto"
	"github.com/thebytearray/BytePayments/internal/database"
	"github.com/thebytearray/BytePayments/model"
	"github.com/thebytearray/BytePayments/repository"
	"github.com/thebytearray/BytePayments/service"
)

// AdminLoginHandler godoc
// @Summary      Admin login
// @Description  Authenticate admin user and return JWT token
// @Tags         admin
// @Accept       json
// @Produce      json
// @Param        request  body  dto.LoginRequest  true  "Login credentials"
// @Success      200  {object}  dto.ApiResponse "Login successful"
// @Failure      400  {object}  dto.ApiResponse "Invalid request"
// @Failure      401  {object}  dto.ApiResponse "Invalid credentials"
// @Router       /api/v1/admin/login [post]
func AdminLoginHandler(ctx *fiber.Ctx) error {
	var req dto.LoginRequest

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(400).JSON(dto.NewError("Invalid request body", err))
	}

	// Validate request
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return ctx.Status(400).JSON(dto.NewError("Validation failed", err))
	}

	adminService := service.NewAdminService(repository.NewAdminRepository(database.DB))

	token, err := adminService.Login(req.Username, req.Password)
	if err != nil {
		return ctx.Status(401).JSON(dto.NewError("Invalid credentials", err))
	}

	// Get admin info for response
	admin, _ := adminService.ValidateToken(token)
	
	response := dto.LoginResponse{
		Token: token,
		Admin: dto.AdminInfo{
			ID:       admin.ID,
			Username: admin.Username,
			Email:    admin.Email,
		},
	}

	return ctx.JSON(dto.NewSuccess("Login successful", response))
}

// Middleware for admin authentication
func AdminAuthMiddleware() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		authHeader := ctx.Get("Authorization")
		if authHeader == "" {
			return ctx.Status(401).JSON(dto.NewError("Authorization header missing", nil))
		}

		// Extract token from "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return ctx.Status(401).JSON(dto.NewError("Invalid authorization header format", nil))
		}

		token := parts[1]
		adminService := service.NewAdminService(repository.NewAdminRepository(database.DB))

		admin, err := adminService.ValidateToken(token)
		if err != nil {
			return ctx.Status(401).JSON(dto.NewError("Invalid or expired token", err))
		}

		// Store admin info in context
		ctx.Locals("admin", admin)
		return ctx.Next()
	}
}

// GetAllPaymentsHandler godoc
// @Summary      Get all payments
// @Description  Get all payments for admin management
// @Tags         admin
// @Accept       json
// @Produce      json
// @Success      200  {object}  dto.ApiResponse "Payments retrieved successfully"
// @Failure      401  {object}  dto.ApiResponse "Unauthorized"
// @Router       /api/v1/admin/payments [get]
func GetAllPaymentsHandler(ctx *fiber.Ctx) error {
	adminService := service.NewAdminManagementService(repository.NewPaymentRepository(database.DB))
	payments, err := adminService.GetAllPayments()
	if err != nil {
		return ctx.Status(500).JSON(dto.NewError("Failed to fetch payments", err))
	}
	return ctx.JSON(dto.NewSuccess("Payments fetched successfully", payments))
}

// DeletePaymentHandler godoc
// @Summary      Delete payment
// @Description  Delete a payment (Admin only)
// @Tags         admin
// @Accept       json
// @Produce      json
// @Param        id  path  string  true  "Payment ID"
// @Success      200  {object}  dto.ApiResponse "Payment deleted successfully"
// @Failure      401  {object}  dto.ApiResponse "Unauthorized"
// @Router       /api/v1/admin/payments/{id} [delete]
func DeletePaymentHandler(ctx *fiber.Ctx) error {
	paymentID := ctx.Params("id")
	if paymentID == "" {
		return ctx.Status(400).JSON(dto.NewError("Payment ID is required", nil))
	}

	adminService := service.NewAdminManagementService(repository.NewPaymentRepository(database.DB))
	err := adminService.DeletePayment(paymentID)
	if err != nil {
		return ctx.Status(500).JSON(dto.NewError("Failed to delete payment", err))
	}
	return ctx.JSON(dto.NewSuccess("Payment deleted successfully", nil))
}

// GetAllWalletsHandler godoc
// @Summary      Get all wallets
// @Description  Get all wallets for admin management
// @Tags         admin
// @Accept       json
// @Produce      json
// @Success      200  {object}  dto.ApiResponse "Wallets retrieved successfully"
// @Failure      401  {object}  dto.ApiResponse "Unauthorized"
// @Router       /api/v1/admin/wallets [get]
func GetAllWalletsHandler(ctx *fiber.Ctx) error {
	adminService := service.NewAdminManagementService(repository.NewPaymentRepository(database.DB))
	wallets, err := adminService.GetAllWallets()
	if err != nil {
		return ctx.Status(500).JSON(dto.NewError("Failed to fetch wallets", err))
	}
	return ctx.JSON(dto.NewSuccess("Wallets fetched successfully", wallets))
}

// DeleteWalletHandler godoc
// @Summary      Delete wallet
// @Description  Delete a wallet (Admin only)
// @Tags         admin
// @Accept       json
// @Produce      json
// @Param        id  path  string  true  "Wallet ID"
// @Success      200  {object}  dto.ApiResponse "Wallet deleted successfully"
// @Failure      401  {object}  dto.ApiResponse "Unauthorized"
// @Router       /api/v1/admin/wallets/{id} [delete]
func DeleteWalletHandler(ctx *fiber.Ctx) error {
	walletID := ctx.Params("id")
	if walletID == "" {
		return ctx.Status(400).JSON(dto.NewError("Wallet ID is required", nil))
	}

	adminService := service.NewAdminManagementService(repository.NewPaymentRepository(database.DB))
	err := adminService.DeleteWallet(walletID)
	if err != nil {
		return ctx.Status(500).JSON(dto.NewError("Failed to delete wallet", err))
	}
	return ctx.JSON(dto.NewSuccess("Wallet deleted successfully", nil))
}

// CreateCurrencyHandler godoc
// @Summary      Create currency
// @Description  Create a new currency (Admin only)
// @Tags         admin
// @Accept       json
// @Produce      json
// @Param        request  body  dto.CreateCurrencyRequest  true  "Currency data"
// @Success      201  {object}  dto.ApiResponse "Currency created successfully"
// @Failure      400  {object}  dto.ApiResponse "Invalid request"
// @Failure      401  {object}  dto.ApiResponse "Unauthorized"
// @Router       /api/v1/admin/currencies [post]
func CreateCurrencyHandler(ctx *fiber.Ctx) error {
	var req dto.CreateCurrencyRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(400).JSON(dto.NewError("Invalid request body", err))
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return ctx.Status(400).JSON(dto.NewError("Validation failed", err))
	}

	currency := &model.Currency{
		Code:         req.Code,
		Name:         req.Name,
		Network:      req.Network,
		IsToken:      req.IsToken,
		ContractAddr: req.ContractAddr,
		Enabled:      req.Enabled,
	}

	currencyService := service.NewCurrenciesService(repository.NewCurrenciesRepository(database.DB))
	err := currencyService.CreateCurrency(currency)
	if err != nil {
		return ctx.Status(400).JSON(dto.NewError("Failed to create currency", err))
	}
	return ctx.Status(201).JSON(dto.NewSuccess("Currency created successfully", currency))
}

// UpdateCurrencyHandler godoc
// @Summary      Update currency
// @Description  Update an existing currency (Admin only)
// @Tags         admin
// @Accept       json
// @Produce      json
// @Param        code     path  string                     true  "Currency Code"
// @Param        request  body  dto.UpdateCurrencyRequest  true  "Currency data"
// @Success      200  {object}  dto.ApiResponse "Currency updated successfully"
// @Failure      400  {object}  dto.ApiResponse "Invalid request"
// @Failure      401  {object}  dto.ApiResponse "Unauthorized"
// @Failure      404  {object}  dto.ApiResponse "Currency not found"
// @Router       /api/v1/admin/currencies/{code} [put]
func UpdateCurrencyHandler(ctx *fiber.Ctx) error {
	currencyCode := ctx.Params("code")
	if currencyCode == "" {
		return ctx.Status(400).JSON(dto.NewError("Currency code is required", nil))
	}

	var req dto.UpdateCurrencyRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(400).JSON(dto.NewError("Invalid request body", err))
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return ctx.Status(400).JSON(dto.NewError("Validation failed", err))
	}

	currencyService := service.NewCurrenciesService(repository.NewCurrenciesRepository(database.DB))

	// Check if currency exists
	existingCurrency, err := currencyService.GetCurrencyByCode(currencyCode)
	if err != nil {
		return ctx.Status(404).JSON(dto.NewError("Currency not found", err))
	}

	// Update currency
	existingCurrency.Code = req.Code
	existingCurrency.Name = req.Name
	existingCurrency.Network = req.Network
	existingCurrency.IsToken = req.IsToken
	existingCurrency.ContractAddr = req.ContractAddr
	existingCurrency.Enabled = req.Enabled

	err = currencyService.UpdateCurrency(existingCurrency)
	if err != nil {
		return ctx.Status(400).JSON(dto.NewError("Failed to update currency", err))
	}
	return ctx.JSON(dto.NewSuccess("Currency updated successfully", existingCurrency))
}

// DeleteCurrencyHandler godoc
// @Summary      Delete currency
// @Description  Delete a currency (Admin only)
// @Tags         admin
// @Accept       json
// @Produce      json
// @Param        code  path  string  true  "Currency Code"
// @Success      200  {object}  dto.ApiResponse "Currency deleted successfully"
// @Failure      401  {object}  dto.ApiResponse "Unauthorized"
// @Failure      404  {object}  dto.ApiResponse "Currency not found"
// @Router       /api/v1/admin/currencies/{code} [delete]
func DeleteCurrencyHandler(ctx *fiber.Ctx) error {
	currencyCode := ctx.Params("code")
	if currencyCode == "" {
		return ctx.Status(400).JSON(dto.NewError("Currency code is required", nil))
	}

	currencyService := service.NewCurrenciesService(repository.NewCurrenciesRepository(database.DB))

	// Check if currency exists
	_, err := currencyService.GetCurrencyByCode(currencyCode)
	if err != nil {
		return ctx.Status(404).JSON(dto.NewError("Currency not found", err))
	}

	err = currencyService.DeleteCurrency(currencyCode)
	if err != nil {
		return ctx.Status(400).JSON(dto.NewError("Failed to delete currency", err))
	}
	return ctx.JSON(dto.NewSuccess("Currency deleted successfully", nil))
}

// ChangePasswordHandler godoc
// @Summary      Change admin password
// @Description  Change admin password with old password verification
// @Tags         admin
// @Accept       json
// @Produce      json
// @Param        request  body  dto.ChangePasswordRequest  true  "Password change data"
// @Success      200  {object}  dto.ApiResponse "Password changed successfully"
// @Failure      400  {object}  dto.ApiResponse "Invalid request"
// @Failure      401  {object}  dto.ApiResponse "Unauthorized or invalid old password"
// @Router       /api/v1/admin/change-password [post]
func ChangePasswordHandler(ctx *fiber.Ctx) error {
	var req dto.ChangePasswordRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(400).JSON(dto.NewError("Invalid request body", err))
	}

	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return ctx.Status(400).JSON(dto.NewError("Validation failed", err))
	}

	// Get admin from context
	admin := ctx.Locals("admin").(*model.Admin)

	adminService := service.NewAdminService(repository.NewAdminRepository(database.DB))
	err := adminService.ChangePassword(admin.ID, req.OldPassword, req.NewPassword)
	if err != nil {
		return ctx.Status(400).JSON(dto.NewError("Failed to change password", err))
	}

	return ctx.JSON(dto.NewSuccess("Password changed successfully", nil))
}