package dto

type CreatePaymentRequest struct {
	PlanId       string `json:"plan_id"`
	Email        string `json:"email"`
	CurrencyCode string `json:"currency_code" validate:"required"`
}

type CreatePaymentResponse struct {
	PaymentId        string  `json:"payment_id"`
	PlanId           string  `json:"plan_id"`
	Email            string  `json:"email"`
	QrImage          string  `json:"qr_image"`
	TrxAmount        float64 `json:"trx_amount"`
	TrxWalletAddress string  `json:"trx_wallet_address"`
}
