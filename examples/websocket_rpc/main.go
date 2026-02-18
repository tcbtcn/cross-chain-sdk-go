package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	wsapi "github.com/tcbtcn/cross-chain-sdk-go/ws-api"
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
		fmt.Println()

		// Once connected, start making RPC calls
		go func() {
			time.Sleep(1 * time.Second) // Wait a bit for connection to stabilize

			// Example 1: Ping/Pong
			fmt.Println("📡 Sending ping...")
			if err := ws.RPC.Ping(); err != nil {
				log.Printf("Failed to send ping: %v", err)
			}

			// Example 2: Get allowed methods
			fmt.Println("📡 Requesting allowed methods...")
			if err := ws.RPC.GetAllowedMethods(); err != nil {
				log.Printf("Failed to get allowed methods: %v", err)
			}

			// Example 3: Get active orders
			fmt.Println("📡 Requesting active orders (page 1, limit 10)...")
			if err := ws.RPC.GetActiveOrders(1, 10); err != nil {
				log.Printf("Failed to get active orders: %v", err)
			}

			// Example 4: Get secrets
			fmt.Println("📡 Requesting secrets (page 1, limit 10)...")
			if err := ws.RPC.GetSecrets(1, 10); err != nil {
				log.Printf("Failed to get secrets: %v", err)
			}
		}()
	})

	ws.OnClose(func() {
		fmt.Println("❌ WebSocket connection closed")
	})

	ws.OnError(func(err error) {
		log.Printf("⚠️  WebSocket error: %v", err)
	})

	// Subscribe to RPC responses

	// Ping/Pong
	ws.RPC.OnPong(func(result string) {
		fmt.Printf("🏓 Pong received: %s\n", result)
		fmt.Println()
	})

	// Get allowed methods
	ws.RPC.OnGetAllowedMethods(func(result []string) {
		fmt.Printf("📋 Allowed Methods:\n")
		for i, method := range result {
			fmt.Printf("   %d. %s\n", i+1, method)
		}
		fmt.Println()
	})

	// Get active orders
	ws.RPC.OnGetActiveOrders(func(result interface{}) {
		fmt.Printf("📦 Active Orders Response:\n")
		resultJSON, err := json.MarshalIndent(result, "   ", "  ")
		if err != nil {
			fmt.Printf("   Error marshaling response: %v\n", err)
		} else {
			fmt.Printf("   %s\n", string(resultJSON))
		}
		fmt.Println()
	})

	// Get secrets
	ws.RPC.OnGetSecrets(func(result interface{}) {
		fmt.Printf("🔑 Secrets Response:\n")
		resultJSON, err := json.MarshalIndent(result, "   ", "  ")
		if err != nil {
			fmt.Printf("   Error marshaling response: %v\n", err)
		} else {
			fmt.Printf("   %s\n", string(resultJSON))
		}
		fmt.Println()
	})

	// Wait for interrupt signal to gracefully shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	fmt.Println("🚀 WebSocket RPC client started.")
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
