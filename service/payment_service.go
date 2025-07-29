package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/thebytearray/BytePayments/dto"
	"github.com/thebytearray/BytePayments/internal/tron"
	"github.com/thebytearray/BytePayments/internal/util"
	"github.com/thebytearray/BytePayments/model"
	"github.com/thebytearray/BytePayments/repository"
	"gorm.io/gorm"
)

type PaymentService interface {
	CreatePayment(ctx context.Context, body dto.CreatePaymentRequest) (dto.PaymentResponse, error)
	CancelPaymentById(id string) dto.ApiResponse
	CheckPaymentStatusById(id string) dto.ApiResponse
}

type paymentService struct {
	repo repository.PaymentRepository
}

func NewPaymentService(repo repository.PaymentRepository) PaymentService {
	return &paymentService{repo}
}

func (s *paymentService) CheckPaymentStatusById(id string) dto.ApiResponse {

	payment, err := s.repo.FindPaymentById(id)

	if err != nil {
		return dto.NewError("Failed to get payment by id.", err)
	}

	base64Image, err := util.GenerateQRCodeBase64(payment.Wallet.WalletAddress)
	if err != nil {
		return dto.NewError("Failed to create qr.", err)
	}

	return dto.NewSuccess("Got payment details successfully.", dto.PaymentResponse{
		PaymentId:        payment.ID,
		Status:           payment.Status,
		PlanId:           payment.PlanID,
		Email:            payment.UserEmail,
		QrImage:          base64Image,
		TrxAmount:        payment.AmountTRX,
		TrxWalletAddress: payment.Wallet.WalletAddress,
	})
}

func (s *paymentService) CancelPaymentById(id string) dto.ApiResponse {
	// Get the payment with id
	payment, err := s.repo.FindPaymentById(id)
	if err != nil {
		return dto.NewError("Failed to find the payment.", err)
	}

	if payment.Status == model.Cancelled {
		return dto.NewSuccess("Payment is already cancelled.", nil)
	}

	if payment.Status == model.Completed {
		return dto.NewError("Payment already completed, can't cancel.", fmt.Errorf("payment cannot be cancelled: %w", errors.New("already completed")))
	}

	// Set status to Cancelled
	payment.Status = model.Cancelled

	// Update in DB
	if err := s.repo.UpdatePayment(payment); err != nil {
		return dto.NewError("Failed to cancel payment", err)
	}

	return dto.NewSuccess("Cancelled payment successfully.", nil)
}

func (s *paymentService) CreatePayment(ctx context.Context, body dto.CreatePaymentRequest) (dto.PaymentResponse, error) {

	//get the selected plan
	plan, err := s.repo.FindPlanById(body.PlanId)

	if err != nil {
		return dto.PaymentResponse{}, fmt.Errorf("plan not found : %w", err)
	}

	//currency
	currency, err := s.repo.FindCurrencyByCode(body.CurrencyCode)

	if err != nil {
		return dto.PaymentResponse{}, fmt.Errorf("curency not found : %w", err)
	}

	//convert usd to trx

	amountTrx, err := tron.ConvertUSDToTRX(plan.PriceUSD)

	if err != nil {
		return dto.PaymentResponse{}, fmt.Errorf("failed to convert amount to trx : %w", err)
	}

	//check if the wallet for user there or not:9

	wallet, err := s.repo.FindWalletByEmail(body.Email)

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {

		return dto.PaymentResponse{}, fmt.Errorf("db error : %w", err)
	}

	var walletID string

	if wallet.ID != "" { //he already has a wallet
		walletID = wallet.ID
	} else {
		//nope create a new for that guy
		//
		walletSecret, walletAddr, err := tron.GenerateWallet()

		if err != nil {
			return dto.PaymentResponse{}, fmt.Errorf("wallet generation error : %w", err)
		}

		encKey, err := util.AesEncryptPK(walletSecret)

		if err != nil {
			return dto.PaymentResponse{}, fmt.Errorf("wallet secret encryption failed : %w", err)
		}

		newWallet := model.Wallet{
			ID:            util.GenerateUniqueID(),
			Email:         body.Email,
			WalletAddress: walletAddr,
			WalletSecret:  encKey,
		}

		err = s.repo.CreateWallet(newWallet)
		if err != nil {
			return dto.PaymentResponse{}, fmt.Errorf("failed to create wallet : %w", err)
		}

		walletID = newWallet.ID
		wallet = newWallet
	}

	payment := model.Payment{
		ID:            util.GenerateUniqueID(),
		PlanID:        plan.ID,
		AmountUSD:     plan.PriceUSD,
		WalletID:      walletID,
		CurrencyCode:  currency.Code,
		AmountTRX:     amountTrx,
		UserEmail:     body.Email,
		Status:        model.Pending,
		PaidAmountTRX: 0,
	}

	err = s.repo.CreatePayment(payment)

	if err != nil {
		return dto.PaymentResponse{}, fmt.Errorf("failed to create payment : %w", err)
	}
	//generate a qr
	base64Image, err := util.GenerateQRCodeBase64(wallet.WalletAddress)
	if err != nil {
		return dto.PaymentResponse{}, fmt.Errorf("failed to create qr : %w", err)
	}
	return dto.PaymentResponse{
		PaymentId:        payment.ID,
		Status:           model.Pending,
		PlanId:           plan.ID,
		Email:            body.Email,
		QrImage:          base64Image,
		TrxAmount:        amountTrx,
		TrxWalletAddress: wallet.WalletAddress,
	}, nil

}
