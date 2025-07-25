package main

import (
	"github.com/TheByteArray/BytePayments/config"
	"github.com/TheByteArray/BytePayments/internal/database"
	"github.com/TheByteArray/BytePayments/route"
	_ "github.com/TheByteArray/BytePayments/docs" // Swagger docs import
)

func main() {
	config.NewConfig()
	database.NewConnection()
	//	database.SeedDatabase()
	app := route.NewRouter()
	app.Listen(":8080")

}

// import (
// 	"crypto/ecdsa"
// 	"encoding/hex"
// 	"fmt"
// 	"log"

// 	"github.com/TheByteArray/go-tron-sdk/pkg/address"
// 	"github.com/TheByteArray/go-tron-sdk/pkg/client"
// 	"github.com/TheByteArray/go-tron-sdk/pkg/client/transaction"
// 	"github.com/TheByteArray/go-tron-sdk/pkg/keys"
// 	"github.com/TheByteArray/go-tron-sdk/pkg/proto/api"
// )

// func main() {
// 	// config.LoadConfig()
// 	// database.Connect()
// 	//
// 	//
// 	//
// 	//
// 	//
// 	//
// 	//c5bb3b9c8ef87c13be48a9c280dda8b8c5bbaaaddad871068b5a3aee19636e36 -> private key
// 	//
// 	//
// 	c := client.NewGrpcClient("grpc.trongrid.io:50051")
// 	err := c.Start(client.GRPCInsecure())
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer c.Stop()
// 	btcecPrivKey, err := keys.GetPrivateKeyFromHex("c5bb3b9c8ef87c13be48a9c280dda8b8c5bbaaaddad871068b5a3aee19636e36")
// 	if err != nil {
// 		// handle error
// 	}

// 	// Convert to *ecdsa.PrivateKey
// 	ecdsaPrivKey := btcecPrivKey.ToECDSA()

// 	// Now pass ecdsaPrivKey to sendTRX

// 	result, err := sendTRX(c, "TJRoFJTWFCV2WbZWDzr7uZkQTYKdr9JGTF", "TM2Z6o6SabAJ3cW8UWjoG3orAGYPcqqdzJ", 10000000000, ecdsaPrivKey)

// 	log.Println(result)
// 	log.Println(err)
// }

// func sendTRX(c *client.GrpcClient, from, to string, amount int64, privateKey *ecdsa.PrivateKey) (string, error) {
// 	// Convert addresses
// 	fromAddr, _ := address.Base58ToAddress(from)
// 	toAddr, _ := address.Base58ToAddress(to)

// 	// Create transaction
// 	tx, err := c.Transfer(fromAddr.String(), toAddr.String(), amount)
// 	if err != nil {
// 		return "", err
// 	}

// 	// Sign transaction
// 	signedTx, err := transaction.SignTransactionECDSA(tx.Transaction, privateKey)
// 	if err != nil {
// 		return "", err
// 	}

// 	// Broadcast
// 	result, err := c.Broadcast(signedTx)
// 	if err != nil {
// 		return "", err
// 	}

// 	if !result.Result ||
// 		result.Code != api.Return_SUCCESS {
// 		return "", fmt.Errorf("transaction failed: (%d) %s", result.Code, result.Message)
// 	}

// 	return hex.EncodeToString(tx.Txid), nil
// }
