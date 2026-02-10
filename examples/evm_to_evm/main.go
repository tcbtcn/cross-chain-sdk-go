package main

import (
	"context"
	"fmt"
	"log"
	"time"

	apiorders "github.com/dawitel/cross-chain-sdk-go/api/orders"
	"github.com/dawitel/cross-chain-sdk-go/chains"
	"github.com/dawitel/cross-chain-sdk-go/domains/hashlock"
	"github.com/dawitel/cross-chain-sdk-go/sdk"
)

func main() {
	config := sdk.Config{
		URL:     "https://api.1inch.dev/fusion-plus",
		AuthKey: "your-auth-key-here",
	}

	client := sdk.NewSDK(config)

	ctx := context.Background()

	// Ethereum to Ethereum (EVM to EVM) swap
	// Example: USDT on Polygon -> BNB on Binance Smart Chain
	quoteParams := sdk.QuoteParams{
		SrcChainID:      chains.Polygon,
		DstChainID:      chains.Binance,
		SrcTokenAddress: "0xc2132d05d31c914a87c6611c10748aeb04b58e8f", // USDT on Polygon
		DstTokenAddress: "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee", // Native BNB
		Amount:          "10000000",                                   // 10 USDT (6 decimals)
		WalletAddress:   "0xYourWalletAddress",
		EnableEstimate:  true,
	}

	quote, err := client.GetQuote(ctx, quoteParams)
	if err != nil {
		log.Fatalf("Failed to get quote: %v", err)
	}

	fmt.Printf("Got quote: %+v\n", quote)

	preset := quote.Presets.Fast
	secrets, err := sdk.GenerateSecrets(preset.SecretsCount)
	if err != nil {
		log.Fatalf("Failed to generate secrets: %v", err)
	}

	secretHashes := make([]string, len(secrets))
	for i, secret := range secrets {
		hash, err := hashlock.HashSecret(secret)
		if err != nil {
			log.Fatalf("Failed to hash secret: %v", err)
		}
		secretHashes[i] = hash
	}

	var hashLock *hashlock.HashLock
	if len(secrets) == 1 {
		hashLock, err = hashlock.ForSingleFill(secrets[0])
		if err != nil {
			log.Fatalf("Failed to create hash lock: %v", err)
		}
	} else {
		leaves, err := hashlock.GetMerkleLeaves(secrets)
		if err != nil {
			log.Fatalf("Failed to get merkle leaves: %v", err)
		}
		hashLock, err = hashlock.ForMultipleFills(leaves)
		if err != nil {
			log.Fatalf("Failed to create hash lock: %v", err)
		}
	}

	orderParams := sdk.OrderParams{
		WalletAddress: quoteParams.WalletAddress,
		HashLock:      hashLock,
		SecretHashes:  secretHashes,
		Preset:        "fast",
	}

	preparedOrder, err := client.CreateOrder(ctx, quote, orderParams)
	if err != nil {
		log.Fatalf("Failed to create order: %v", err)
	}

	fmt.Printf("Created order with hash: %s\n", preparedOrder.Hash)

	for {
		readyToAccept, err := client.GetReadyToAcceptSecretFills(ctx, preparedOrder.Hash)
		if err != nil {
			log.Fatalf("Failed to get ready to accept secret fills: %v", err)
		}

		if len(readyToAccept.Fills) > 0 {
			for _, fill := range readyToAccept.Fills {
				err := client.SubmitSecret(ctx, preparedOrder.Hash, secrets[fill.Idx])
				if err != nil {
					log.Printf("Failed to submit secret: %v", err)
					continue
				}
				fmt.Printf("Submitted secret for fill %d\n", fill.Idx)
			}
		}

		status, err := client.GetOrderStatus(ctx, preparedOrder.Hash)
		if err != nil {
			log.Fatalf("Failed to get order status: %v", err)
		}

		fmt.Printf("Order status: %s\n", status.Status)

		if status.Status == apiorders.OrderStatusExecuted ||
			status.Status == apiorders.OrderStatusExpired ||
			status.Status == apiorders.OrderStatusRefunded {
			break
		}

		time.Sleep(5 * time.Second)
	}

	fmt.Println("Order completed")
}
