package controller

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/thebytearray/BytePayments/dto"
	"github.com/thebytearray/BytePayments/internal/database"
	"github.com/thebytearray/BytePayments/internal/tron"
	"github.com/thebytearray/BytePayments/internal/util"
	"github.com/thebytearray/BytePayments/model"
)

func CreatePaymentHandler(ctx *fiber.Ctx) error {
	//get the email,and plan id
	var body dto.CreatePaymentRequest

	err := ctx.BodyParser(&body)

	validate := validator.New()
	if err := validate.Struct(&body); err != nil {
		return ctx.JSON(dto.ApiResponse{
			Status:     string(dto.ERROR),
			StatusCode: http.StatusBadRequest,
			Error:      string(err.Error()),
			Message:    fmt.Sprintf("Invalid request body, %v", err),
			Data:       nil,
		})
	}

	if err != nil {
		log.Println(err)
		return ctx.JSON(dto.ApiResponse{
			Status:     string(dto.ERROR),
			StatusCode: http.StatusBadRequest,
			Error:      string(err.Error()),
			Message:    fmt.Sprintf("Invalid request body, %v", err),
			Data:       nil,
		})
	}
	// check if user has a wallet in db or not
	var userWallet model.Wallet
	walletResult := database.DB.Where("email = ?", body.Email).First(&userWallet)
	//fetch the plan price in trx
	var selectedPlan model.Plan
	planResult := database.DB.Where("id = ?", body.PlanId).First(&selectedPlan)

	//fetch the currency and validate
	//

	var selectedCurrency model.Currency

	currencyResult := database.DB.Where("code = ?", body.CurrencyCode).First(&selectedCurrency)

	if currencyResult.RowsAffected == 0 || currencyResult.Error != nil {
		return ctx.JSON(dto.ApiResponse{
			Status:     string(dto.ERROR),
			StatusCode: http.StatusNotFound,
			Error:      string(currencyResult.Error.Error()),
			Message:    fmt.Sprintf("Currency not found, error : %v", currencyResult.Error),
			Data:       nil,
		})
	}

	if planResult.RowsAffected == 0 || planResult.Error != nil {
		return ctx.JSON(dto.ApiResponse{
			Status:     string(dto.ERROR),
			StatusCode: http.StatusNotFound,
			Error:      string(planResult.Error.Error()),
			Message:    fmt.Sprintf("Plan not found, error : %v", planResult.Error),
			Data:       nil,
		})
	}
	// convert the price to trx
	amountTrx, err := tron.ConvertUSDToTRX(selectedPlan.PriceUSD)

	if err != nil {
		return ctx.JSON(dto.ApiResponse{
			Status:     string(dto.ERROR),
			StatusCode: http.StatusNotFound,
			Error:      string(err.Error()),
			Message:    fmt.Sprintf("Failed to convert price in trx, error : %v", err.Error()),
			Data:       nil,
		})
	}

	if walletResult.RowsAffected != 0 {
		//we have the wallet just response it.

		//	create the payment in db
		newPayment := model.Payment{
			ID:            uuid.NewString(),
			PlanID:        selectedPlan.ID,
			AmountUSD:     selectedPlan.PriceUSD,
			WalletID:      userWallet.ID,
			CurrencyCode:  body.CurrencyCode,
			AmountTRX:     amountTrx,
			UserEmail:     body.Email,
			Status:        model.Pending,
			PaidAmountTRX: float64(0),
		}

		createPaymentResult := database.DB.Create(&newPayment)

		if createPaymentResult.Error != nil || createPaymentResult.RowsAffected == 0 {
			return ctx.JSON(dto.ApiResponse{
				Status:     string(dto.ERROR),
				StatusCode: http.StatusExpectationFailed,
				Error:      string(err.Error()),
				Message:    fmt.Sprintf("Failed to create payment, error : %v", createPaymentResult.Error.Error()),
				Data:       nil,
			})
		}

		return ctx.JSON(dto.ApiResponse{
			Status:     string(dto.OK),
			StatusCode: http.StatusOK,
			Error:      "nil",
			Message:    "Payment created successfully.",
			Data: dto.CreatePaymentResponse{
				PaymentId:        newPayment.ID,
				PlanId:           newPayment.PlanID,
				Email:            newPayment.UserEmail,
				TrxAmount:        amountTrx,
				TrxWalletAddress: userWallet.WalletAddress,
			},
		})

	}

	// show if have create a new payment with that details, and tell to pay to that address,
	// if dont have then create a new wallet and create a new payment and tell to pay in that address
	// //create a new wallet and insert with the email

	walletSecret, walletAddr, err := tron.GenerateWallet()

	if err != nil {
		return ctx.JSON(dto.ApiResponse{
			Status:     "error",
			StatusCode: http.StatusExpectationFailed,
			Error:      string(err.Error()),
			Message:    fmt.Sprintf("Failed to create wallet, error : %v", err.Error()),
			Data:       nil,
		})
	}

	encryptedKey, err := util.AesEncryptPK(walletSecret)

	if err != nil {
		return ctx.JSON(dto.ApiResponse{
			Status:     "error",
			StatusCode: http.StatusExpectationFailed,
			Error:      string(err.Error()),
			Message:    fmt.Sprintf("Failed to encrypt wallet, error : %v", err.Error()),
			Data:       nil,
		})
	}

	//create a wallet struct with data
	newWallet := model.Wallet{
		ID:            uuid.NewString(),
		Email:         body.Email,
		WalletAddress: walletAddr,
		WalletSecret:  encryptedKey,
	}

	newWalletResult := database.DB.Create(&newWallet)

	if newWalletResult.RowsAffected == 0 || newWalletResult.Error != nil {
		return ctx.JSON(dto.ApiResponse{
			Status:     "error",
			StatusCode: http.StatusExpectationFailed,
			Error:      string(newWalletResult.Error.Error()),
			Message:    fmt.Sprintf("Failed to create wallet, error : %v", newWalletResult.Error.Error()),
			Data:       nil,
		})
	}

	newWalletPayment := model.Payment{
		ID:            uuid.NewString(),
		PlanID:        selectedPlan.ID,
		AmountUSD:     selectedPlan.PriceUSD,
		WalletID:      newWallet.ID,
		CurrencyCode:  body.CurrencyCode,
		AmountTRX:     amountTrx,
		UserEmail:     body.Email,
		Status:        model.Pending,
		PaidAmountTRX: float64(0),
	}

	createNewPaymentResult := database.DB.Create(&newWalletPayment)

	if createNewPaymentResult.Error != nil || createNewPaymentResult.RowsAffected == 0 {
		return ctx.JSON(dto.ApiResponse{
			Status:     "error",
			StatusCode: http.StatusExpectationFailed,
			Error:      string(createNewPaymentResult.Error.Error()),
			Message:    fmt.Sprintf("Failed to create payment, error : %v", createNewPaymentResult.Error.Error()),
			Data:       nil,
		})
	}

	return ctx.JSON(dto.ApiResponse{
		Status:     "ok",
		StatusCode: http.StatusOK,
		Error:      "nil",
		Message:    "Payment created successfully.",
		Data: dto.CreatePaymentResponse{
			PaymentId:        newWalletPayment.ID,
			PlanId:           newWalletPayment.PlanID,
			Email:            newWalletPayment.UserEmail,
			TrxAmount:        amountTrx,
			TrxWalletAddress: newWallet.WalletAddress,
		},
	})
}

func CancelPaymentHandler(ctx *fiber.Ctx) error {

	return nil
}

func GetPaymentStatusHandler(ctx *fiber.Ctx) error {

	return nil
}
