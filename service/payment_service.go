package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/thebytearray/BytePayments/dto"
	"github.com/thebytearray/BytePayments/internal/tron"
	"github.com/thebytearray/BytePayments/internal/util"
	"github.com/thebytearray/BytePayments/model"
	"github.com/thebytearray/BytePayments/repository"
	"gorm.io/gorm"
)

type PaymentService interface {
	CreatePayment(ctx context.Context, body dto.CreatePaymentRequest) (dto.CreatePaymentResponse, error)
}

type paymentService struct {
	repo repository.PaymentRepository
}

func NewPaymentService(repo repository.PaymentRepository) PaymentService {
	return &paymentService{repo}
}

func (s *paymentService) CreatePayment(ctx context.Context, body dto.CreatePaymentRequest) (dto.CreatePaymentResponse, error) {

	//get the selected plan
	plan, err := s.repo.FindPlanById(body.PlanId)

	if err != nil {
		return dto.CreatePaymentResponse{}, fmt.Errorf("plan not found : %w", err)
	}

	//currency
	currency, err := s.repo.FindCurrencyByCode(body.CurrencyCode)

	if err != nil {
		return dto.CreatePaymentResponse{}, fmt.Errorf("curency not found : %w", err)
	}

	//convert usd to trx

	amountTrx, err := tron.ConvertUSDToTRX(plan.PriceUSD)

	if err != nil {
		return dto.CreatePaymentResponse{}, fmt.Errorf("failed to convert amount to trx : %w", err)
	}

	//check if the wallet for user there or not:9

	wallet, err := s.repo.FindWalletByEmail(body.Email)

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {

		return dto.CreatePaymentResponse{}, fmt.Errorf("db error : %w", err)
	}

	var walletID string

	if wallet.ID != "" { //he already has a wallet
		walletID = wallet.ID
	} else {
		//nope create a new for that guy
		//
		walletSecret, walletAddr, err := tron.GenerateWallet()

		if err != nil {
			return dto.CreatePaymentResponse{}, fmt.Errorf("wallet generation error : %w", err)
		}

		encKey, err := util.AesEncryptPK(walletSecret)

		if err != nil {
			return dto.CreatePaymentResponse{}, fmt.Errorf("wallet secret encryption failed : %w", err)
		}

		newWallet := model.Wallet{

			ID:            uuid.NewString(),
			Email:         body.Email,
			WalletAddress: walletAddr,
			WalletSecret:  encKey,
		}

		err = s.repo.CreateWallet(newWallet)
		if err != nil {
			return dto.CreatePaymentResponse{}, fmt.Errorf("failed to create wallet : %w", err)
		}

		walletID = newWallet.ID
		wallet = newWallet
	}

	payment := model.Payment{
		ID:            uuid.NewString(),
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
		return dto.CreatePaymentResponse{}, fmt.Errorf("failed to create payment : %w", err)
	}

	return dto.CreatePaymentResponse{
		PaymentId:        payment.ID,
		PlanId:           plan.ID,
		Email:            body.Email,
		TrxAmount:        amountTrx,
		TrxWalletAddress: wallet.WalletAddress,
	}, nil

}
