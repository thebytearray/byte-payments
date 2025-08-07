package service

import (
	"github.com/thebytearray/BytePayments/model"
	"github.com/thebytearray/BytePayments/repository"
)

type AdminManagementService interface {
	// Payment Management
	GetAllPayments() ([]model.Payment, error)
	DeletePayment(id string) error
	// Wallet Management
	GetAllWallets() ([]model.Wallet, error)
	DeleteWallet(id string) error
}

type adminManagementService struct {
	paymentRepo repository.PaymentRepository
}

func NewAdminManagementService(paymentRepo repository.PaymentRepository) AdminManagementService {
	return &adminManagementService{
		paymentRepo: paymentRepo,
	}
}

func (s *adminManagementService) GetAllPayments() ([]model.Payment, error) {
	return s.paymentRepo.GetAllPayments()
}

func (s *adminManagementService) DeletePayment(id string) error {
	return s.paymentRepo.DeletePayment(id)
}

func (s *adminManagementService) GetAllWallets() ([]model.Wallet, error) {
	return s.paymentRepo.GetAllWallets()
}

func (s *adminManagementService) DeleteWallet(id string) error {
	return s.paymentRepo.DeleteWallet(id)
}