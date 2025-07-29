package service

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"github.com/dgraph-io/ristretto"
	"github.com/thebytearray/BytePayments/internal/database"
	"github.com/thebytearray/BytePayments/internal/util"
)

type VerificationService interface {
	GenerateAndSendCode(email string) error
	VerifyCode(email, code string) (bool, string, error)
	IsEmailVerified(verificationToken string) bool
}

type verificationService struct {
	emailService EmailService
	cache        *ristretto.Cache
}

func NewVerificationService() VerificationService {
	return &verificationService{
		emailService: NewEmailService(),
		cache:        database.Cache,
	}
}

// GenerateAndSendCode generates a 6-digit code and sends it via email
func (v *verificationService) GenerateAndSendCode(email string) error {
	// Generate 6-digit code
	code, err := v.generateCode()
	if err != nil {
		return fmt.Errorf("failed to generate verification code: %w", err)
	}

	// Store code in cache with 10 minutes expiration
	codeKey := fmt.Sprintf("verification_code:%s", email)
	v.cache.SetWithTTL(codeKey, code, 1, 10*time.Minute)

	// Send email
	err = v.emailService.SendVerificationCode(email, code)
	if err != nil {
		// Clean up cache if email sending fails
		v.cache.Del(codeKey)
		return fmt.Errorf("failed to send verification email: %w", err)
	}

	return nil
}

// VerifyCode checks if the provided code matches the stored code
func (v *verificationService) VerifyCode(email, code string) (bool, string, error) {
	codeKey := fmt.Sprintf("verification_code:%s", email)
	
	// Get stored code
	storedCodeValue, found := v.cache.Get(codeKey)
	if !found {
		return false, "", fmt.Errorf("verification code not found or expired")
	}

	storedCode, ok := storedCodeValue.(string)
	if !ok {
		return false, "", fmt.Errorf("invalid stored code format")
	}

	// Check if codes match
	if storedCode != code {
		return false, "", nil
	}

	// Code is valid, generate verification token
	verificationToken := util.GenerateUniqueID()
	tokenKey := fmt.Sprintf("verified_email:%s", verificationToken)
	
	// Store verification token for 30 minutes
	v.cache.SetWithTTL(tokenKey, email, 1, 30*time.Minute)

	// Delete the used code
	v.cache.Del(codeKey)

	return true, verificationToken, nil
}

// IsEmailVerified checks if the verification token is valid
func (v *verificationService) IsEmailVerified(verificationToken string) bool {
	tokenKey := fmt.Sprintf("verified_email:%s", verificationToken)
	
	_, found := v.cache.Get(tokenKey)
	return found
}

// generateCode creates a random 6-digit code
func (v *verificationService) generateCode() (string, error) {
	code := ""
	for i := 0; i < 6; i++ {
		digit, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", err
		}
		code += digit.String()
	}
	return code, nil
} 