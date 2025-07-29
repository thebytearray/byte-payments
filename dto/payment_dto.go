package dto

import "github.com/thebytearray/BytePayments/model"

type CreatePaymentRequest struct {
	PlanId            string `json:"plan_id"`
	Email             string `json:"email"`
	CurrencyCode      string `json:"currency_code" validate:"required"`
	VerificationToken string `json:"verification_token" validate:"required"`
}

type PaymentResponse struct {
	PaymentId        string              `json:"payment_id"`
	Status           model.PaymentStatus `json:"status"`
	PlanId           string              `json:"plan_id"`
	Email            string              `json:"email"`
	QrImage          string              `json:"qr_image"`
	TrxAmount        float64             `json:"trx_amount"`
	TrxWalletAddress string              `json:"trx_wallet_address"`
	CreatedAt        string              `json:"created_at"`
	UpdatedAt        string              `json:"updated_at"`
}
