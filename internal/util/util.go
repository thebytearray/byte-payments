package util

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"

	"github.com/skip2/go-qrcode"

	"github.com/thebytearray/BytePayments/config"
)

func AesEncryptPK(privateKey string) (string, error) {
	key := []byte(config.Cfg.TRX_WALLET_ENCRYPTION_KEY)
	if len(key) != 32 {
		return "", fmt.Errorf("encryption key must be 32 bytes (got %d)", len(key))
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	nonce := make([]byte, aesGCM.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	ciphertext := aesGCM.Seal(nil, nonce, []byte(privateKey), nil)

	// prepend nonce to ciphertext
	final := append(nonce, ciphertext...)

	// return as base64 encoded string
	return base64.StdEncoding.EncodeToString(final), nil
}

func AesDecryptPK(encrypted string) (string, error) {
	key := []byte(config.Cfg.TRX_WALLET_ENCRYPTION_KEY)
	if len(key) != 32 {
		return "", fmt.Errorf("decryption key must be 32 bytes")
	}

	data, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return "", fmt.Errorf("base64 decode failed: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	if len(data) < aesGCM.NonceSize() {
		return "", fmt.Errorf("ciphertext too short")
	}

	nonce := data[:aesGCM.NonceSize()]
	ciphertext := data[aesGCM.NonceSize():]

	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("decryption failed: %w", err)
	}

	return string(plaintext), nil
}

func GenerateQRCodeBase64(content string) (string, error) {
	png, err := qrcode.Encode(content, qrcode.Medium, 256)
	if err != nil {
		return "", fmt.Errorf("failed to encode QR: %w", err)
	}

	base64Image := base64.StdEncoding.EncodeToString(png)
	return "data:image/png;base64," + base64Image, nil
}
