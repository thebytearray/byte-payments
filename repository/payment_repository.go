package repository

import (
	"log"
	"time"

	"github.com/thebytearray/BytePayments/model"
	"gorm.io/gorm"
)

type PaymentRepository interface {
	FindWalletByEmail(email string) (model.Wallet, error)
	FindPlanById(id string) (model.Plan, error)
	FindCurrencyByCode(code string) (model.Currency, error)
	CreateWallet(wallet model.Wallet) error
	CreatePayment(payment model.Payment) error
	FindPaymentById(id string) (model.Payment, error)
	UpdatePayment(payment model.Payment) error
	HasPendingPayment(user_email string) (bool, error)
	FindAllPendingPayments() ([]model.Payment, error)
	MarkAsCompletedById(id string, paidAmount float64, completedAt *time.Time) error
	MarkAsExpiredById(id string) error
	// Admin methods
	GetAllPayments() ([]model.Payment, error)
	DeletePayment(id string) error
	GetAllWallets() ([]model.Wallet, error)
	DeleteWallet(id string) error
}

type paymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) PaymentRepository {
	return &paymentRepository{db}
}

func (r *paymentRepository) FindAllPendingPayments() ([]model.Payment, error) {
	var payments []model.Payment
	err := r.db.Where("status = ? AND created_at >= ?", model.Pending, time.Now().Add(-15*time.Minute)).
		Preload("Wallet").
		Preload("Plan").
		Find(&payments).Error
	return payments, err
}

func (r *paymentRepository) MarkAsCompletedById(id string, paidAmount float64, completedAt *time.Time) error {
	return r.db.Model(&model.Payment{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"status":          model.Completed,
			"paid_amount_trx": paidAmount,
			"updated_at":      completedAt,
		}).Error
}

func (r *paymentRepository) MarkAsExpiredById(id string) error {
	return r.db.Model(&model.Payment{}).Where("id = ?", id).Update("status", model.Expired).Error
}

func (r *paymentRepository) HasPendingPayment(user_email string) (bool, error) {
	var count int64
	err := r.db.Model(&model.Payment{}).Where("user_email = ? AND status = ?", user_email, model.Pending).Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, err
}

func (r *paymentRepository) UpdatePayment(payment model.Payment) error {
	return r.db.Save(&payment).Error
}

func (r *paymentRepository) FindWalletByEmail(email string) (model.Wallet, error) {
	var wallet model.Wallet
	res := r.db.Where("email = ?", email).First(&wallet)
	return wallet, res.Error
}

func (r *paymentRepository) FindPlanById(id string) (model.Plan, error) {
	var plan model.Plan
	res := r.db.Where("id = ?", id).First(&plan)

	return plan, res.Error
}

func (r *paymentRepository) FindCurrencyByCode(code string) (model.Currency, error) {
	var currency model.Currency
	res := r.db.Where("code = ?", code).First(&currency)

	return currency, res.Error
}

func (r *paymentRepository) CreateWallet(wallet model.Wallet) error {
	return r.db.Create(&wallet).Error
}

func (r *paymentRepository) CreatePayment(payment model.Payment) error {
	return r.db.Create(&payment).Error
}

func (r *paymentRepository) FindPaymentById(id string) (model.Payment, error) {
	var payment model.Payment
	res := r.db.Preload("Wallet").Where("id = ?", id).Find(&payment)
	log.Println(payment.CurrencyCode)
	return payment, res.Error
}

// Admin methods
func (r *paymentRepository) GetAllPayments() ([]model.Payment, error) {
	var payments []model.Payment
	res := r.db.Preload("Wallet").Preload("Plan").Order("created_at DESC").Find(&payments)
	return payments, res.Error
}

func (r *paymentRepository) DeletePayment(id string) error {
	return r.db.Delete(&model.Payment{}, "id = ?", id).Error
}

func (r *paymentRepository) GetAllWallets() ([]model.Wallet, error) {
	var wallets []model.Wallet
	res := r.db.Order("created_at DESC").Find(&wallets)
	return wallets, res.Error
}

func (r *paymentRepository) DeleteWallet(id string) error {
	return r.db.Delete(&model.Wallet{}, "id = ?", id).Error
}
