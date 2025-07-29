package service

import (
	"errors"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/thebytearray/BytePayments/config"
	"github.com/thebytearray/BytePayments/dto"
	"github.com/thebytearray/BytePayments/internal/tron"
	"github.com/thebytearray/BytePayments/internal/util"
	"github.com/thebytearray/BytePayments/model"
	"github.com/thebytearray/BytePayments/repository"
	"gorm.io/gorm"
)

type PaymentService interface {
	CreatePayment(body dto.CreatePaymentRequest) (dto.PaymentResponse, error)
	CancelPaymentById(id string) dto.ApiResponse
	CheckPaymentStatusById(id string) dto.ApiResponse
	ProcessPendingPayments()
}

type paymentService struct {
	repo repository.PaymentRepository
}

func NewPaymentService(repo repository.PaymentRepository) PaymentService {
	return &paymentService{repo}
}

func (s *paymentService) ProcessPendingPayments() {
	const completionThreshold = 0.95 // 95%
	const tolerance = 0.001
	payments, err := s.repo.FindAllPendingPayments()
	if err != nil {
		log.Println("Error fetching pending payments : ", err)
		return
	}

	emailService := NewEmailService()

	for _, p := range payments {
		//check for expiry
		if time.Since(p.CreatedAt) > 15*time.Minute {
			err := s.repo.MarkAsExpiredById(p.ID)
			if err != nil {
				log.Println("failed to mark payment as expired", err)
			}
			continue
		}

		//check actual balance?
		balance, err := tron.CheckBalance(tron.TRON_CLIENT, p.Wallet.WalletAddress)

		if err != nil {
			log.Println("Error fetching balance : ", err)
			continue
		}

		//parse
		diff := balance - p.AmountTRX
		absDiff := math.Abs(diff)
		receivedRatio := balance / p.AmountTRX

		if receivedRatio >= completionThreshold || absDiff <= tolerance {
			// Payment amount is sufficient, but don't mark as completed yet
			log.Printf("Payment %s has sufficient funds: received %.6f TRX (expected %.6f)", p.ID, balance, p.AmountTRX)
			
			// Try to sweep funds first
			err = s.sweepFunds(p)
			if err != nil {
				log.Printf("Failed to sweep funds for payment %s: %v", p.ID, err)
				// Don't mark as completed, will retry on next cron run
				continue
			}
			
			// Sweeping successful, now mark as completed
			now := time.Now()
			err = s.repo.MarkAsCompletedById(p.ID, balance, &now)
			if err != nil {
				log.Printf("Failed to mark completed for payment %s: %v", p.ID, err)
				continue
			}

			// Update the payment object with the paid amount for email templates
			p.Status = model.Completed
			p.PaidAmountTRX = balance
			p.UpdatedAt = now

			log.Printf("Payment %s completed and funds swept successfully", p.ID)
			
			// Determine payment condition and send appropriate email
			if balance > p.AmountTRX+tolerance {
				// Overpaid
				overpaidAmount := balance - p.AmountTRX
				log.Printf("Payment %s overpaid by %.6f TRX", p.ID, overpaidAmount)
				
				err = emailService.SendOverpaymentEmail(p, p.Plan, overpaidAmount)
				if err != nil {
					log.Printf("Failed to send overpayment email for payment %s: %v", p.ID, err)
				} else {
					log.Printf("Overpayment email sent for payment %s", p.ID)
				}
			} else {
				// Exact or close enough payment
				err = emailService.SendPaymentCompletionEmail(p, p.Plan)
				if err != nil {
					log.Printf("Failed to send completion email for payment %s: %v", p.ID, err)
				} else {
					log.Printf("Completion email sent for payment %s", p.ID)
				}
			}
		} else if balance > 0 && balance < p.AmountTRX-tolerance {
			// Underpaid - check if we haven't already sent an email recently
			remainingAmount := p.AmountTRX - balance
			log.Printf("Payment %s underpaid: received %.6f TRX, remaining %.6f TRX", p.ID, balance, remainingAmount)
			
			// Only send underpayment email if we have received some payment and haven't sent one recently
			// You might want to add a field to track when the last email was sent to avoid spam
			if balance > 0.001 { // Only if they've paid something significant
				// Create a temporary payment object with the current balance for email
				tempPayment := p
				tempPayment.PaidAmountTRX = balance
				
				err = emailService.SendUnderpaymentEmail(tempPayment, p.Plan, remainingAmount)
				if err != nil {
					log.Printf("Failed to send underpayment email for payment %s: %v", p.ID, err)
				} else {
					log.Printf("Underpayment email sent for payment %s", p.ID)
				}
			}
		} else {
			log.Printf("Payment %s still pending: received %.6f TRX (%.2f%% of expected)", p.ID, balance, receivedRatio*100)
		}

	}

}

func (s *paymentService) sweepFunds(payment model.Payment) error {
	// 1. Check wallet balance
	balance, err := tron.CheckBalance(tron.TRON_CLIENT, payment.Wallet.WalletAddress)
	if err != nil {
		return fmt.Errorf("failed to check balance: %w", err)
	}

	transferable, err := tron.GetTransferableAmount(payment.Wallet.WalletAddress, balance)
	if err != nil {
		return fmt.Errorf("failed to calculate transferable amount: %w", err)
	}

	if transferable <= 0 {
		return fmt.Errorf("no transferable amount available")
	}

	mainWalletAddr := config.Cfg.TRX_HOT_WALLET_ADDRESS
	paymentWalletPrivKey, err := util.AesDecryptPK(payment.Wallet.WalletSecret)

	if err != nil {
		return fmt.Errorf("failed to decrypt wallet key: %w", err)
	}

	txID, err := tron.SendTRX(tron.TRON_CLIENT, payment.Wallet.WalletAddress, mainWalletAddr, transferable, paymentWalletPrivKey)
	if err != nil {
		return fmt.Errorf("failed to send trx: %w", err)
	}

	log.Printf("Swept %.6f TRX from %s to main wallet. TxID: %s", transferable, payment.Wallet.WalletAddress, txID)
	return nil
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
		CreatedAt:        payment.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:        payment.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
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

func (s *paymentService) CreatePayment(body dto.CreatePaymentRequest) (dto.PaymentResponse, error) {
	// Verify email verification token first
	verificationService := NewVerificationService()
	if !verificationService.IsEmailVerified(body.VerificationToken) {
		return dto.PaymentResponse{}, fmt.Errorf("email verification required. Please verify your email first")
	}

	hasPendingPayment, err := s.repo.HasPendingPayment(body.Email)
	if err != nil {
		return dto.PaymentResponse{}, fmt.Errorf("failed to check pending payment : %w", err)
	}

	if hasPendingPayment {
		return dto.PaymentResponse{}, fmt.Errorf("please complete or cancel the previous pending payment to create a new one.")
	}
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
		CreatedAt:        payment.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:        payment.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil

}
