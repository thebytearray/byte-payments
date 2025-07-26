package tron

import (
	"encoding/hex"
	"fmt"

	"github.com/TheByteArray/go-tron-sdk/pkg/address"
	"github.com/TheByteArray/go-tron-sdk/pkg/client"
	"github.com/TheByteArray/go-tron-sdk/pkg/client/transaction"
	"github.com/TheByteArray/go-tron-sdk/pkg/keys"
	"github.com/TheByteArray/go-tron-sdk/pkg/proto/api"
)

func CheckBalance(c *client.GrpcClient, addr string) (string, error) {
	tronAddr, err := address.Base58ToAddress(addr)
	if err != nil {
		return "", fmt.Errorf("invalid TRON address: %w", err)
	}

	account, err := c.GetAccount(tronAddr.String())
	if err != nil {
		return "", fmt.Errorf("failed to get account: %w", err)
	}

	balanceTRX := float64(account.Balance) / 1e6
	return fmt.Sprintf("%.6f TRX", balanceTRX), nil
}

// GenerateWallet creates a new TRON wallet and returns private key hex and base58 address.
func GenerateWallet() (privKeyHex string, base58Addr string, err error) {
	privateKey, err := keys.GenerateKey()
	if err != nil {
		return "", "", fmt.Errorf("failed to generate private key: %w", err)
	}

	privKeyHex = fmt.Sprintf("%x", privateKey.Serialize())
	addrBytes := address.BTCECPrivkeyToAddress(privateKey)
	base58Addr = address.Address(addrBytes).String()

	return privKeyHex, base58Addr, nil
}
func SendTRX(c *client.GrpcClient, from, to string, amountTRX float64, privateKey string) (string, error) {
	// Convert TRX to sun (1 TRX = 1,000,000 sun)
	amountSun := int64(amountTRX * 1_000_000)

	// Parse the private key
	btcecPrivKey, err := keys.GetPrivateKeyFromHex(privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to parse private key: %w", err)
	}
	ecdsaPrivKey := btcecPrivKey.ToECDSA()

	// Convert from and to addresses
	fromAddr, err := address.Base58ToAddress(from)
	if err != nil {
		return "", fmt.Errorf("invalid from address: %w", err)
	}
	toAddr, err := address.Base58ToAddress(to)
	if err != nil {
		return "", fmt.Errorf("invalid to address: %w", err)
	}

	// Create the transfer transaction
	tx, err := c.Transfer(fromAddr.String(), toAddr.String(), amountSun)
	if err != nil {
		return "", fmt.Errorf("failed to create transfer transaction: %w", err)
	}

	// Sign the transaction
	signedTx, err := transaction.SignTransactionECDSA(tx.Transaction, ecdsaPrivKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign transaction: %w", err)
	}

	// Broadcast the transaction
	result, err := c.Broadcast(signedTx)
	if err != nil {
		return "", fmt.Errorf("failed to broadcast transaction: %w", err)
	}
	if !result.Result || result.Code != api.Return_SUCCESS {
		return "", fmt.Errorf("transaction rejected by network: (%d) %s", result.Code, result.Message)
	}

	// Return transaction ID
	return hex.EncodeToString(tx.Txid), nil
}
