package tron

import (
	"log"

	"github.com/TheByteArray/go-tron-sdk/pkg/client"
	"github.com/thebytearray/BytePayments/config"
)

var TRON_CLIENT *client.GrpcClient

func NewClient() {
	var grpcURL string

	if config.Cfg.APP_ENV == "development" {
		grpcURL = config.Cfg.TRON_GRPC_TESTNET
	} else {
		grpcURL = config.Cfg.TRON_GRPC_MAINNET
	}

	c := client.NewGrpcClient(grpcURL)

	if err := c.Start(client.GRPCInsecure()); err != nil {
		log.Fatalf("failed to start TRON client: %v", err)
	}

	TRON_CLIENT = c
}
