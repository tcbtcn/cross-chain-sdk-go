# Cross-Chain SDK for Go

Production-grade Go SDK for creating atomic cross-chain token swaps through 1inch Fusion Plus API.

## Installation

```bash
go get github.com/dawitel/cross-chain-sdk-go
```

## Quick Start

```go
package main

import (
    "context"
    "log"

    "github.com/dawitel/cross-chain-sdk-go/chains"
    "github.com/dawitel/cross-chain-sdk-go/domains/hashlock"
    "github.com/dawitel/cross-chain-sdk-go/sdk"
)

func main() {
    config := sdk.Config{
        URL:     "https://api.1inch.dev/fusion-plus",
        AuthKey: "your-auth-key",
    }

    client := sdk.NewSDK(config)
    ctx := context.Background()

    quote, err := client.GetQuote(ctx, sdk.QuoteParams{
        SrcChainID:      chains.Polygon,
        DstChainID:      chains.Binance,
        SrcTokenAddress: "0xc2132d05d31c914a87c6611c10748aeb04b58e8f",
        DstTokenAddress: "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee",
        Amount:          "10000000",
        WalletAddress:   "0xYourWalletAddress",
        EnableEstimate:  true,
    })
    if err != nil {
        log.Fatal(err)
    }

    secrets, _ := sdk.GenerateSecrets(quote.Presets.Fast.SecretsCount)
    secretHashes := make([]string, len(secrets))
    for i, secret := range secrets {
        hash, _ := hashlock.HashSecret(secret)
        secretHashes[i] = hash
    }

    hashLock, _ := hashlock.ForSingleFill(secrets[0])

    order, err := client.CreateOrder(ctx, quote, sdk.OrderParams{
        WalletAddress: "0xYourWalletAddress",
        HashLock:      hashLock,
        SecretHashes:  secretHashes,
        Preset:        "fast",
    })
    if err != nil {
        log.Fatal(err)
    }

    log.Printf("Order created: %s", order.Hash)
}
```

## Features

- Cross-chain token swaps (EVM to EVM, EVM to Solana, Solana to EVM)
- Quote generation with multiple presets (fast, medium, slow)
- Order creation and submission
- Secret management for escrow deployments
- Order status tracking
- Support for native assets
- EIP-712 signing for EVM orders
- **WebSocket API** for real-time order updates and RPC functionality

## Supported Chains

- Ethereum
- Polygon
- Binance Smart Chain
- Arbitrum
- Avalanche
- Optimism
- Gnosis
- Coinbase
- zkSync
- Linea
- Sonic
- Unichain
- Solana

## API Reference

### SDK Client

#### NewSDK

Creates a new SDK client instance.

```go
client := sdk.NewSDK(sdk.Config{
    URL:     "https://api.1inch.dev/fusion-plus",
    AuthKey: "your-auth-key",
    BlockchainProvider: yourProvider, // optional
})
```

#### GetQuote

Gets a quote for a cross-chain swap.

```go
quote, err := client.GetQuote(ctx, sdk.QuoteParams{
    SrcChainID:      chains.Polygon,
    DstChainID:      chains.Binance,
    SrcTokenAddress: "0x...",
    DstTokenAddress: "0x...",
    Amount:          "10000000",
    WalletAddress:   "0x...",
    EnableEstimate:  true,
})
```

#### CreateOrder

Creates an order from a quote.

```go
order, err := client.CreateOrder(ctx, quote, sdk.OrderParams{
    WalletAddress: "0x...",
    HashLock:      hashLock,
    SecretHashes:  secretHashes,
    Preset:        "fast",
})
```

#### SubmitOrder

Submits an EVM order to the relayer.

```go
orderInfo, err := client.SubmitOrder(ctx, chains.Polygon, evmOrder, quoteID, secretHashes)
```

#### GetOrderStatus

Gets the status of an order.

```go
status, err := client.GetOrderStatus(ctx, orderHash)
```

#### GetReadyToAcceptSecretFills

Gets fills that are ready to accept secrets.

```go
ready, err := client.GetReadyToAcceptSecretFills(ctx, orderHash)
```

#### SubmitSecret

Submits a secret for an escrow deployment.

```go
err := client.SubmitSecret(ctx, orderHash, secret)
```

### WebSocket API

The WebSocket API provides real-time order updates and RPC functionality for the cross-chain SDK.

#### NewWebSocketApi

Creates a new WebSocket API client.

```go
import "github.com/dawitel/cross-chain-sdk-go/ws-api"

ws, err := wsapi.NewWebSocketApi(wsapi.WsApiConfig{
    URL:     "wss://api.1inch.dev/fusion-plus/ws",
    AuthKey: "your-auth-key",
    LazyInit: false, // Set to true for lazy initialization
})
if err != nil {
    log.Fatal(err)
}

// Initialize connection (if lazyInit was true)
ws.Init()
```

#### Order Events

Subscribe to real-time order updates.

```go
// Subscribe to all order events
ws.Order.OnOrder(func(data wsapi.OrderEventType) {
    switch data.GetEvent() {
    case wsapi.EventTypeOrderCreated:
        if created, ok := data.(*wsapi.OrderCreatedEvent); ok {
            fmt.Printf("Order created: %s\n", created.OrderHash)
        }
    case wsapi.EventTypeOrderFilled:
        if filled, ok := data.(*wsapi.OrderFilledEvent); ok {
            fmt.Printf("Order filled: %s\n", filled.OrderHash)
        }
    }
})

// Subscribe to specific order events
ws.Order.OnOrderCreated(func(data *wsapi.OrderCreatedEvent) {
    fmt.Printf("Order created: %s\n", data.OrderHash)
})

ws.Order.OnOrderFilled(func(data *wsapi.OrderFilledEvent) {
    fmt.Printf("Order filled: %s\n", data.OrderHash)
})

ws.Order.OnOrderCancelled(func(data *wsapi.OrderCancelledEvent) {
    fmt.Printf("Order cancelled: %s\n", data.OrderHash)
})

ws.Order.OnOrderInvalid(func(data *wsapi.OrderInvalidEvent) {
    fmt.Printf("Order invalid: %s, reason: %s\n", data.OrderHash, data.Reason)
})
```

#### RPC Methods

Call RPC methods for querying data.

```go
// Ping/Pong
ws.RPC.Ping()
ws.RPC.OnPong(func(result string) {
    fmt.Printf("Pong received: %s\n", result)
})

// Get active orders
ws.RPC.GetActiveOrders(1, 10) // page, limit
ws.RPC.OnGetActiveOrders(func(result interface{}) {
    fmt.Printf("Active orders received\n")
})

// Get secrets
ws.RPC.GetSecrets(1, 10) // page, limit
ws.RPC.OnGetSecrets(func(result interface{}) {
    fmt.Printf("Secrets received\n")
})

// Get allowed methods
ws.RPC.GetAllowedMethods()
ws.RPC.OnGetAllowedMethods(func(result []string) {
    fmt.Printf("Allowed methods: %v\n", result)
})
```

#### Connection Management

Handle WebSocket connection lifecycle.

```go
// Connection events
ws.OnOpen(func() {
    fmt.Println("WebSocket connected")
})

ws.OnClose(func() {
    fmt.Println("WebSocket disconnected")
})

ws.OnError(func(err error) {
    fmt.Printf("WebSocket error: %v\n", err)
})

// Close connection
defer ws.Close()
```

#### WebSocket Features

- **Order Events**: Subscribe to real-time order updates (created, filled, cancelled, invalid, etc.)
- **RPC Methods**: Call RPC methods like `getActiveOrders`, `getSecrets`, `ping`
- **Connection Management**: Handle connection lifecycle (open, close, error)
- **Type Safety**: Strongly typed event structures matching the TypeScript SDK

## Examples

See the `examples/` directory for complete working examples:

- `examples/ethereum_to_ethereum/` - Ethereum to Ethereum (EVM to EVM) swap example
- `examples/evm_to_evm/` - EVM to EVM cross-chain swap example (Polygon to BSC)
- `examples/evm_to_solana/` - EVM to Solana cross-chain swap example (Ethereum to Solana)
- `examples/solana_to_evm/` - Solana to EVM cross-chain swap example (Solana to Ethereum)
- `examples/websocket_order_events/` - WebSocket API example for subscribing to order events
- `examples/websocket_rpc/` - WebSocket API example for RPC methods

## Requirements

- Go 1.21 or higher
- Valid 1inch API auth key from [1inch Portal](https://portal.1inch.dev)

## License

See LICENSE file for details.
