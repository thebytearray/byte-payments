package service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/thebytearray/BytePayments/config"
	"github.com/thebytearray/BytePayments/model"
	"github.com/thebytearray/BytePayments/repository"
	"golang.org/x/crypto/bcrypt"
)

type AdminService interface {
	Login(username, password string) (string, error)
	CreateAdmin(username, email, password string) error
	ValidateToken(tokenString string) (*model.Admin, error)
	ChangePassword(adminID uint, oldPassword, newPassword string) error
}

type adminService struct {
	repo repository.AdminRepository
}

func NewAdminService(repo repository.AdminRepository) AdminService {
	return &adminService{repo}
}

func (s *adminService) Login(username, password string) (string, error) {
	admin, err := s.repo.GetAdminByUsername(username)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	// Check password
	err = bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"admin_id": admin.ID,
		"username": admin.Username,
		"email":    admin.Email,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // 24 hour expiry
	})

	tokenString, err := token.SignedString([]byte(config.Cfg.JWT_SECRET))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *adminService) CreateAdmin(username, email, password string) error {
	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	admin := &model.Admin{
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
		IsActive: true,
	}

	return s.repo.CreateAdmin(admin)
}

func (s *adminService) ValidateToken(tokenString string) (*model.Admin, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Cfg.JWT_SECRET), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		username := claims["username"].(string)
		return s.repo.GetAdminByUsername(username)
	}

	return nil, errors.New("invalid token")
}

func (s *adminService) ChangePassword(adminID uint, oldPassword, newPassword string) error {
	// Get admin by ID
	admin, err := s.repo.GetAdminByID(adminID)
	if err != nil {
		return errors.New("admin not found")
	}

	// Verify old password
	err = bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(oldPassword))
	if err != nil {
		return errors.New("invalid old password")
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Update password
	admin.Password = string(hashedPassword)
	return s.repo.UpdateAdmin(admin)
}