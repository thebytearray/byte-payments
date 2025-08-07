package controller

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/thebytearray/BytePayments/dto"
	"github.com/thebytearray/BytePayments/internal/database"
	"github.com/thebytearray/BytePayments/internal/util"
	"github.com/thebytearray/BytePayments/model"
	"github.com/thebytearray/BytePayments/repository"
	"github.com/thebytearray/BytePayments/service"
)

// GetPlansHandler godoc
// @Summary      Get all subscription plans
// @Description  Returns a list of all available subscription plans with their pricing and features
// @Tags         plans
// @Accept       json
// @Produce      json
// @Success      200  {object}  dto.ApiResponse "Plans retrieved successfully"
// @Failure      404  {object}  dto.ApiResponse "No plans found"
// @Router       /api/v1/plans [get]
func GetPlansHandler(ctx *fiber.Ctx) error {

	plansService := service.NewPlansService(repository.NewPlansRepository(database.DB))

	plans, err := plansService.GetPlans()

	if err != nil {
		return ctx.JSON(dto.NewError("Plan not found", err))
	}

	return ctx.JSON(dto.NewSuccess("Plans fetched successfully.", plans))
}

// CreatePlanHandler godoc
// @Summary      Create a new plan
// @Description  Create a new subscription plan (Admin only)
// @Tags         admin
// @Accept       json
// @Produce      json
// @Param        request  body  dto.CreatePlanRequest  true  "Plan data"
// @Success      201  {object}  dto.ApiResponse "Plan created successfully"
// @Failure      400  {object}  dto.ApiResponse "Invalid request"
// @Failure      401  {object}  dto.ApiResponse "Unauthorized"
// @Router       /api/v1/admin/plans [post]
func CreatePlanHandler(ctx *fiber.Ctx) error {
	var req dto.CreatePlanRequest

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(400).JSON(dto.NewError("Invalid request body", err))
	}

	// Validate request
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return ctx.Status(400).JSON(dto.NewError("Validation failed", err))
	}

	plansService := service.NewPlansService(repository.NewPlansRepository(database.DB))

	plan := &model.Plan{
		ID:           util.GenerateUniqueID(),
		Name:         req.Name,
		Description:  req.Description,
		PriceUSD:     req.PriceUSD,
		DurationDays: req.DurationDays,
	}

	err := plansService.CreatePlan(plan)
	if err != nil {
		return ctx.Status(400).JSON(dto.NewError("Failed to create plan", err))
	}

	return ctx.Status(201).JSON(dto.NewSuccess("Plan created successfully", plan))
}

// UpdatePlanHandler godoc
// @Summary      Update a plan
// @Description  Update an existing subscription plan (Admin only)
// @Tags         admin
// @Accept       json
// @Produce      json
// @Param        id       path  string                 true  "Plan ID"
// @Param        request  body  dto.UpdatePlanRequest  true  "Plan data"
// @Success      200  {object}  dto.ApiResponse "Plan updated successfully"
// @Failure      400  {object}  dto.ApiResponse "Invalid request"
// @Failure      401  {object}  dto.ApiResponse "Unauthorized"
// @Failure      404  {object}  dto.ApiResponse "Plan not found"
// @Router       /api/v1/admin/plans/{id} [put]
func UpdatePlanHandler(ctx *fiber.Ctx) error {
	planID := ctx.Params("id")
	if planID == "" {
		return ctx.Status(400).JSON(dto.NewError("Plan ID is required", nil))
	}

	var req dto.UpdatePlanRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(400).JSON(dto.NewError("Invalid request body", err))
	}

	// Validate request
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		return ctx.Status(400).JSON(dto.NewError("Validation failed", err))
	}

	plansService := service.NewPlansService(repository.NewPlansRepository(database.DB))

	// Check if plan exists
	existingPlan, err := plansService.GetPlanByID(planID)
	if err != nil {
		return ctx.Status(404).JSON(dto.NewError("Plan not found", err))
	}

	// Update plan
	existingPlan.Name = req.Name
	existingPlan.Description = req.Description
	existingPlan.PriceUSD = req.PriceUSD
	existingPlan.DurationDays = req.DurationDays

	err = plansService.UpdatePlan(existingPlan)
	if err != nil {
		return ctx.Status(400).JSON(dto.NewError("Failed to update plan", err))
	}

	return ctx.JSON(dto.NewSuccess("Plan updated successfully", existingPlan))
}

// DeletePlanHandler godoc
// @Summary      Delete a plan
// @Description  Delete a subscription plan (Admin only)
// @Tags         admin
// @Accept       json
// @Produce      json
// @Param        id  path  string  true  "Plan ID"
// @Success      200  {object}  dto.ApiResponse "Plan deleted successfully"
// @Failure      401  {object}  dto.ApiResponse "Unauthorized"
// @Failure      404  {object}  dto.ApiResponse "Plan not found"
// @Router       /api/v1/admin/plans/{id} [delete]
func DeletePlanHandler(ctx *fiber.Ctx) error {
	planID := ctx.Params("id")
	if planID == "" {
		return ctx.Status(400).JSON(dto.NewError("Plan ID is required", nil))
	}

	plansService := service.NewPlansService(repository.NewPlansRepository(database.DB))

	// Check if plan exists
	_, err := plansService.GetPlanByID(planID)
	if err != nil {
		return ctx.Status(404).JSON(dto.NewError("Plan not found", err))
	}

	err = plansService.DeletePlan(planID)
	if err != nil {
		return ctx.Status(400).JSON(dto.NewError("Failed to delete plan", err))
	}

	return ctx.JSON(dto.NewSuccess("Plan deleted successfully", nil))
}
