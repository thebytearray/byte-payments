package repository

import (
	"github.com/thebytearray/BytePayments/model"
	"gorm.io/gorm"
)

type AdminRepository interface {
	GetAdminByUsername(username string) (*model.Admin, error)
	CreateAdmin(admin *model.Admin) error
	GetAdminByID(id uint) (*model.Admin, error)
	UpdateAdmin(admin *model.Admin) error
}

type adminRepository struct {
	db *gorm.DB
}

func NewAdminRepository(db *gorm.DB) AdminRepository {
	return &adminRepository{db}
}

func (r *adminRepository) GetAdminByUsername(username string) (*model.Admin, error) {
	var admin model.Admin
	res := r.db.Where("username = ? AND is_active = ?", username, true).First(&admin)
	if res.Error != nil {
		return nil, res.Error
	}
	return &admin, nil
}

func (r *adminRepository) CreateAdmin(admin *model.Admin) error {
	return r.db.Create(admin).Error
}

func (r *adminRepository) GetAdminByID(id uint) (*model.Admin, error) {
	var admin model.Admin
	res := r.db.Where("id = ? AND is_active = ?", id, true).First(&admin)
	if res.Error != nil {
		return nil, res.Error
	}
	return &admin, nil
}

func (r *adminRepository) UpdateAdmin(admin *model.Admin) error {
	return r.db.Save(admin).Error
}