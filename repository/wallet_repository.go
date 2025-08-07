package repository

import (
	"github.com/thebytearray/BytePayments/model"
	"gorm.io/gorm"
)

type WalletRepository interface {
	GetWallets() ([]model.Wallet, error)
	GetWalletByID(id uint) (*model.Wallet, error)
	CreateWallet(wallet *model.Wallet) error
	UpdateWallet(wallet *model.Wallet) error
	DeleteWallet(id uint) error
}

type walletRepository struct {
	db *gorm.DB
}

func NewWalletRepository(db *gorm.DB) WalletRepository {
	return &walletRepository{db}
}

func (r *walletRepository) GetWallets() ([]model.Wallet, error) {
	var wallets []model.Wallet
	res := r.db.Find(&wallets)
	return wallets, res.Error
}

func (r *walletRepository) GetWalletByID(id uint) (*model.Wallet, error) {
	var wallet model.Wallet
	res := r.db.First(&wallet, id)
	if res.Error != nil {
		return nil, res.Error
	}
	return &wallet, nil
}

func (r *walletRepository) CreateWallet(wallet *model.Wallet) error {
	return r.db.Create(wallet).Error
}

func (r *walletRepository) UpdateWallet(wallet *model.Wallet) error {
	return r.db.Save(wallet).Error
}

func (r *walletRepository) DeleteWallet(id uint) error {
	return r.db.Delete(&model.Wallet{}, id).Error
}