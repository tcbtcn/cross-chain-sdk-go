package main

import (
	"context"
	"fmt"
	"log"
	"time"

	apiorders "github.com/tcbtcn/cross-chain-sdk-go/api/orders"
	"github.com/tcbtcn/cross-chain-sdk-go/chains"
	"github.com/tcbtcn/cross-chain-sdk-go/domains/hashlock"
	"github.com/tcbtcn/cross-chain-sdk-go/orders"
	"github.com/tcbtcn/cross-chain-sdk-go/sdk"
)

func main() {
	config := sdk.Config{
		URL:     "https://api.1inch.dev/fusion-plus",
		AuthKey: "your-auth-key-here",
	}

	client := sdk.NewSDK(config)

	ctx := context.Background()

	// Solana to Ethereum (EVM) swap
	// Example: USDC on Solana -> USDC on Ethereum
	quoteParams := sdk.QuoteParams{
		SrcChainID:      chains.Solana,
		DstChainID:      chains.Ethereum,
		SrcTokenAddress: "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGxrHJf99", // USDC on Solana
		DstTokenAddress: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48", // USDC on Ethereum
		Amount:          "100000000",                                  // 100 USDC (6 decimals)
		WalletAddress:   "YourSolanaWalletAddress",
		EnableEstimate:  true,
	}

	quote, err := client.GetQuote(ctx, quoteParams)
	if err != nil {
		log.Fatalf("Failed to get quote: %v", err)
	}

	fmt.Printf("Got quote: QuoteID=%s, SrcAmount=%s, DstAmount=%s\n",
		quote.QuoteID, quote.SrcTokenAmount, quote.DstTokenAmount)

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

	// For Solana orders, receiver is required and must be an EVM address
	orderParams := sdk.OrderParams{
		WalletAddress: quoteParams.WalletAddress,
		HashLock:      hashLock,
		SecretHashes:  secretHashes,
		Preset:        "fast",
		Receiver:      "0xYourEthereumReceiverAddress", // Required for Solana orders
	}

	preparedOrder, err := client.CreateOrder(ctx, quote, orderParams)
	if err != nil {
		log.Fatalf("Failed to create order: %v", err)
	}

	fmt.Printf("Created order with hash: %s\n", preparedOrder.Hash)

	// For Solana orders, we need to announce first before on-chain creation
	solanaOrder, ok := preparedOrder.Order.(*orders.SolanaCrossChainOrder)
	if !ok {
		log.Fatalf("Expected Solana order, got different type")
	}

	// Announce order to relayer (required before on-chain creation)
	orderHash, err := client.AnnounceOrder(ctx, solanaOrder, preparedOrder.QuoteID, secretHashes)
	if err != nil {
		log.Fatalf("Failed to announce order: %v", err)
	}

	fmt.Printf("Announced Solana order with hash: %s\n", orderHash)
	fmt.Println("NOTE: After announcing, you need to create the order on-chain using Solana program")
	fmt.Println("The order must be created on Solana blockchain before it can be filled")

	// Monitor order status and submit secrets when ready
	for {
		readyToAccept, err := client.GetReadyToAcceptSecretFills(ctx, orderHash)
		if err != nil {
			log.Fatalf("Failed to get ready to accept secret fills: %v", err)
		}

		if len(readyToAccept.Fills) > 0 {
			for _, fill := range readyToAccept.Fills {
				err := client.SubmitSecret(ctx, orderHash, secrets[fill.Idx])
				if err != nil {
					log.Printf("Failed to submit secret: %v", err)
					continue
				}
				fmt.Printf("Submitted secret for fill %d\n", fill.Idx)
			}
		}

		status, err := client.GetOrderStatus(ctx, orderHash)
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
