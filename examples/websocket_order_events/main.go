package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dawitel/cross-chain-sdk-go/ws-api"
)

func main() {
	// Get auth key from environment or use placeholder
	authKey := os.Getenv("INCH_AUTH_KEY")
	if authKey == "" {
		authKey = "your-auth-key-here"
		log.Println("Warning: Using placeholder auth key. Set INCH_AUTH_KEY environment variable.")
	}

	// Create WebSocket API client
	ws, err := wsapi.NewWebSocketApi(wsapi.WsApiConfig{
		URL:      "wss://api.1inch.dev/fusion-plus/ws",
		AuthKey:  authKey,
		LazyInit: false, // Auto-connect on creation
	})
	if err != nil {
		log.Fatalf("Failed to create WebSocket client: %v", err)
	}

	// Handle connection events
	ws.OnOpen(func() {
		fmt.Println("✅ WebSocket connected successfully")
	})

	ws.OnClose(func() {
		fmt.Println("❌ WebSocket connection closed")
	})

	ws.OnError(func(err error) {
		log.Printf("⚠️  WebSocket error: %v", err)
	})

	// Subscribe to all order events
	ws.Order.OnOrder(func(data wsapi.OrderEventType) {
		switch data.GetEvent() {
		case wsapi.EventTypeOrderCreated:
			if created, ok := data.(*wsapi.OrderCreatedEvent); ok {
				fmt.Printf("📝 Order Created: %s (Src: %d, Dst: %d)\n",
					created.OrderHash, created.SrcChainID, created.DstChainID)
			}

		case wsapi.EventTypeOrderFilled:
			if filled, ok := data.(*wsapi.OrderFilledEvent); ok {
				fmt.Printf("✅ Order Filled: %s\n", filled.OrderHash)
			}

		case wsapi.EventTypeOrderFilledPartially:
			if partial, ok := data.(*wsapi.OrderFilledPartiallyEvent); ok {
				fmt.Printf("🔄 Order Partially Filled: %s\n", partial.OrderHash)
			}

		case wsapi.EventTypeOrderCancelled:
			if cancelled, ok := data.(*wsapi.OrderCancelledEvent); ok {
				fmt.Printf("❌ Order Cancelled: %s (Remaining: %s)\n",
					cancelled.OrderHash, cancelled.RemainingMakerAmount)
			}

		case wsapi.EventTypeOrderInvalid:
			if invalid, ok := data.(*wsapi.OrderInvalidEvent); ok {
				fmt.Printf("⚠️  Order Invalid: %s", invalid.OrderHash)
				if invalid.Reason != "" {
					fmt.Printf(" - Reason: %s", invalid.Reason)
				}
				fmt.Println()
			}

		case wsapi.EventTypeOrderBalanceChange:
			if balance, ok := data.(*wsapi.OrderBalanceChangeEvent); ok {
				fmt.Printf("💰 Order Balance Change: %s (Remaining: %s, Balance: %s)\n",
					balance.OrderHash, balance.RemainingMakerAmount, balance.Balance)
			}

		case wsapi.EventTypeOrderAllowanceChange:
			if allowance, ok := data.(*wsapi.OrderAllowanceChangeEvent); ok {
				fmt.Printf("🔐 Order Allowance Change: %s (Remaining: %s, Allowance: %s)\n",
					allowance.OrderHash, allowance.RemainingMakerAmount, allowance.Allowance)
			}

		case wsapi.EventTypeOrderSecretShared:
			if secret, ok := data.(*wsapi.OrderSecretSharedEvent); ok {
				fmt.Printf("🔑 Secret Shared: Index %d for order\n", secret.Idx)
			}
		}
	})

	// Also subscribe to specific events for more detailed handling
	ws.Order.OnOrderCreated(func(data *wsapi.OrderCreatedEvent) {
		fmt.Printf("📋 Order Details:\n")
		fmt.Printf("   Hash: %s\n", data.OrderHash)
		fmt.Printf("   Quote ID: %s\n", data.QuoteID)
		fmt.Printf("   Chains: %d -> %d\n", data.SrcChainID, data.DstChainID)
		fmt.Printf("   Maker: %s\n", data.Order.Maker)
		fmt.Printf("   Making Amount: %s\n", data.Order.MakingAmount)
		fmt.Printf("   Taking Amount: %s\n", data.Order.TakingAmount)
		if len(data.SecretHashes) > 0 {
			fmt.Printf("   Secret Hashes: %d\n", len(data.SecretHashes))
		}
		fmt.Println()
	})

	ws.Order.OnOrderFilled(func(data *wsapi.OrderFilledEvent) {
		fmt.Printf("🎉 Order Successfully Filled: %s\n", data.OrderHash)
		if data.TxHash != "" {
			fmt.Printf("   Transaction Hash: %s\n", data.TxHash)
		}
		fmt.Println()
	})

	// Wait for interrupt signal to gracefully shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	fmt.Println("🚀 WebSocket client started. Listening for order events...")
	fmt.Println("   Press Ctrl+C to exit")
	fmt.Println()

	// Keep the connection alive
	<-sigChan
	fmt.Println("\n🛑 Shutting down...")

	// Close the WebSocket connection
	if err := ws.Close(); err != nil {
		log.Printf("Error closing WebSocket: %v", err)
	}

	// Give a moment for cleanup
	time.Sleep(500 * time.Millisecond)
	fmt.Println("✅ Shutdown complete")
}
