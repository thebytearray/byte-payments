package tron

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/TheByteArray/go-tron-sdk/pkg/address"
	"github.com/TheByteArray/go-tron-sdk/pkg/client"
	"github.com/TheByteArray/go-tron-sdk/pkg/client/transaction"
	"github.com/TheByteArray/go-tron-sdk/pkg/keys"
	"github.com/TheByteArray/go-tron-sdk/pkg/proto/api"
	"github.com/thebytearray/BytePayments/config"
	"github.com/thebytearray/BytePayments/dto"
)

var (
	ErrInvalidAddress      = errors.New("invalid TRON address")
	ErrInsufficientBalance = errors.New("insufficient balance to cover transaction fee")
)

func CheckBalance(c *client.GrpcClient, addr string) (float64, error) {
	tronAddr, err := address.Base58ToAddress(addr)
	if err != nil {
		return float64(0), fmt.Errorf("invalid TRON address: %w", err)
	}

	account, err := c.GetAccount(tronAddr.String())
	if err != nil {
		return float64(0), fmt.Errorf("failed to get account: %w", err)
	}

	balanceTRX := float64(account.Balance) / 1e6
	return balanceTRX, nil
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

func ConvertUSDToTRX(usdAmount float64) (float64, error) {
	// Fetch current TRX/USDT price from Binance
	resp, err := http.Get(config.Cfg.BINANCE_API_URL)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch price: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("non-200 response from Binance: %d", resp.StatusCode)
	}

	var priceResp dto.PriceResponse
	if err := json.NewDecoder(resp.Body).Decode(&priceResp); err != nil {
		return 0, fmt.Errorf("failed to decode response: %w", err)
	}

	// Convert string price to float
	price, err := strconv.ParseFloat(priceResp.Price, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid price format: %w", err)
	}

	// Calculate and return TRX amount
	trxAmount := usdAmount / price
	return trxAmount, nil
}

func GetTransferableAmount(walletAddress string, balanceTRX float64) (float64, error) {
	if !strings.HasPrefix(walletAddress, "T") {
		return 0, ErrInvalidAddress
	}

	apiKey := config.Cfg.TRON_GRID_API_KEY
	env := strings.ToLower(config.Cfg.APP_ENV)

	var tronGridURL string
	switch env {
	case "production":
		tronGridURL = config.Cfg.TRON_GRID_API_URL_MAINNET
	case "development", "staging", "test":
		tronGridURL = config.Cfg.TRON_GRID_API_URL_TESTNET
	default:
		return 0, fmt.Errorf("unsupported APP_ENV: %s", config.Cfg.APP_ENV)
	}

	reqBody := dto.AccountResourceRequest{Address: walletAddress}
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return 0, fmt.Errorf("failed to encode request: %w", err)
	}

	req, err := http.NewRequest("POST", tronGridURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return 0, fmt.Errorf("failed to build request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("TRON-PRO-API-KEY", apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("API request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return 0, fmt.Errorf("TronGrid error (%d): %s", resp.StatusCode, string(bodyBytes))
	}

	var res dto.AccountResourceResponse
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return 0, fmt.Errorf("failed to decode response: %w", err)
	}

	balanceSun := int64(balanceTRX * 1_000_000)

	const (
		trxTxSizeBytes      = int64(300)
		sunPerBandwidthUnit = int64(1000)
	)

	freeBandwidth := res.FreeNetLimit - res.FreeNetUsed
	if freeBandwidth < 0 {
		freeBandwidth = 0
	}

	chargeableBandwidth := trxTxSizeBytes - freeBandwidth
	if chargeableBandwidth < 0 {
		chargeableBandwidth = 0
	}

	// Calculate bandwidth fee
	bandwidthFeeSun := chargeableBandwidth * sunPerBandwidthUnit

	// Energy fee estimation (for simple TRX transfer, energy usage is minimal but add safety)
	energyFeeSun := int64(10_000) // ~0.01 TRX for potential energy costs

	// Total estimated fee
	totalFeeSun := bandwidthFeeSun + energyFeeSun

	// Increased safety buffer significantly to prevent insufficient balance errors
	const safetyBuffer = int64(50_000)

	totalCostSun := totalFeeSun + safetyBuffer

	if balanceSun <= totalCostSun {
		return 0, ErrInsufficientBalance
	}

	transferableSun := balanceSun - totalCostSun

	if transferableSun <= 0 {
		return 0, ErrInsufficientBalance
	}

	// Convert back to TRX with proper precision handling
	transferableTRX := float64(transferableSun) / 1_000_000

	// Use ceiling instead of floor to ensure we don't try to send more than calculated
	// Round down to 6 decimal places to match TRX precision
	return math.Floor(transferableTRX*1e6) / 1e6, nil
}

func TrxToSun(trx float64) int64 {
	return int64(trx * 1_000_000)
}

// SunToTrx converts sun to TRX (with rounding to 6 decimal places)
func SunToTrx(sun int64) float64 {
	trx := float64(sun) / 1_000_000
	return math.Floor(trx*1e6) / 1e6
}
