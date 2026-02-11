package sdk

import (
	"context"
	"encoding/json"
	"math/big"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	clienthttp "github.com/dawitel/cross-chain-sdk-go/api/http"
	apiorders "github.com/dawitel/cross-chain-sdk-go/api/orders"
	"github.com/dawitel/cross-chain-sdk-go/api/quoter"
	"github.com/dawitel/cross-chain-sdk-go/chains"
	"github.com/dawitel/cross-chain-sdk-go/domains/addresses"
	"github.com/dawitel/cross-chain-sdk-go/domains/auction"
	"github.com/dawitel/cross-chain-sdk-go/domains/hashlock"
	"github.com/dawitel/cross-chain-sdk-go/orders"
	"github.com/dawitel/cross-chain-sdk-go/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func sampleQuote() *quoter.QuoterResponse {
	return &quoter.QuoterResponse{
		QuoteID:         "test-quote-id-123",
		SrcChainID:      chains.Ethereum,
		DstChainID:      chains.Polygon,
		SrcTokenAddress: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48",
		DstTokenAddress: "0x2791bca1f2de4661ed88a30c99a7a9449aa84174",
		SrcTokenAmount:  "100000000",
		DstTokenAmount:  "99500000",
		Presets: quoter.QuoterPresets{
			Fast: quoter.PresetData{
				AuctionDuration:    300,
				StartAuctionIn:     0,
				InitialRateBump:    100,
				AuctionStartAmount: "100000000",
				StartAmount:        "100000000",
				AuctionEndAmount:   "99500000",
				CostInDstToken:     "500000",
				Points:             []auction.AuctionPoint{},
				AllowPartialFills:  true,
				AllowMultipleFills: true,
				GasCost: quoter.GasCost{
					GasBumpEstimate:  100000,
					GasPriceEstimate: "2000000000",
				},
				ExclusiveResolver: "0x0000000000000000000000000000000000000000",
				SecretsCount:      1,
			},
			Medium: quoter.PresetData{
				AuctionDuration:    600,
				SecretsCount:       1,
				AllowMultipleFills: true,
				AllowPartialFills:  true,
			},
			Slow: quoter.PresetData{
				AuctionDuration:    900,
				SecretsCount:       1,
				AllowMultipleFills: true,
				AllowPartialFills:  true,
			},
		},
		RecommendedPreset: quoter.PresetFast,
		SrcEscrowFactory:  "0x1111111111111111111111111111111111111111",
		DstEscrowFactory:  "0x2222222222222222222222222222222222222222",
		TimeLocks: quoter.TimeLocksRaw{
			SrcWithdrawal:         3600,
			SrcPublicWithdrawal:   7200,
			SrcCancellation:       10800,
			SrcPublicCancellation: 14400,
			DstWithdrawal:         3600,
			DstPublicWithdrawal:   7200,
			DstCancellation:       10800,
		},
		SrcSafetyDeposit: "1000000",
		DstSafetyDeposit: "1000000",
		AutoK:            1,
	}
}

func sampleEvmOrder() *orders.EvmCrossChainOrder {
	// Use NewEvmAddress which doesn't return an error for valid addresses
	maker := addresses.NewEvmAddress("0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb")
	makerAsset := addresses.NewEvmAddress("0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48")
	takerAsset := addresses.NewEvmAddress("0x2791bca1f2de4661ed88a30c99a7a9449aa84174")
	receiver := maker

	secret := repeatHex("00", 31) + "01"
	hashLock, _ := hashlock.ForSingleFill(secret)

	return &orders.EvmCrossChainOrder{
		Maker:                maker,
		MakerAsset:           makerAsset,
		TakerAsset:           takerAsset,
		MakingAmount:         big.NewInt(100000000),
		TakingAmount:         big.NewInt(99500000),
		Receiver:             receiver,
		Deadline:             big.NewInt(1735689600),
		AuctionStartTime:     big.NewInt(1735689300),
		AuctionEndTime:       big.NewInt(1735689600),
		Nonce:                big.NewInt(12345),
		Salt:                 big.NewInt(67890),
		HashLock:             hashLock,
		DstChainID:           chains.Polygon,
		Extension:            "",
		MultipleFillsAllowed: true,
	}
}

func repeatHex(s string, n int) string {
	return "0x" + strings.Repeat(s, n)
}

func TestNewSDK(t *testing.T) {
	config := Config{
		URL:     "https://api.example.com",
		AuthKey: "test-key",
	}

	sdk := NewSDK(config)
	assert.NotNil(t, sdk)
	assert.NotNil(t, sdk.quoterApi)
	assert.NotNil(t, sdk.ordersApi)
	assert.NotNil(t, sdk.relayerApi)
}

func TestNewSDK_WithCustomHTTPClient(t *testing.T) {
	customClient := clienthttp.NewHTTPClient("https://custom.com", "custom-key")
	config := Config{
		URL:        "https://api.example.com",
		AuthKey:    "test-key",
		HTTPClient: customClient,
	}

	sdk := NewSDK(config)
	assert.NotNil(t, sdk)
}

func TestSDK_GetQuote(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		_ = json.NewEncoder(w).Encode(sampleQuote())
	}))
	defer server.Close()

	config := Config{
		URL:     server.URL,
		AuthKey: "test-key",
	}
	sdk := NewSDK(config)

	params := QuoteParams{
		SrcChainID:      chains.Ethereum,
		DstChainID:      chains.Polygon,
		SrcTokenAddress: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48",
		DstTokenAddress: "0x2791bca1f2de4661ed88a30c99a7a9449aa84174",
		Amount:          "100000000",
		WalletAddress:   "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb",
		EnableEstimate:  true,
	}

	quote, err := sdk.GetQuote(context.Background(), params)
	require.NoError(t, err)
	assert.NotNil(t, quote)
	assert.Equal(t, params.SrcChainID, quote.SrcChainID)
}

func TestSDK_GetQuoteWithCustomPreset(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		_ = json.NewEncoder(w).Encode(sampleQuote())
	}))
	defer server.Close()

	config := Config{
		URL:     server.URL,
		AuthKey: "test-key",
	}
	sdk := NewSDK(config)

	params := QuoteParams{
		SrcChainID:      chains.Ethereum,
		DstChainID:      chains.Polygon,
		SrcTokenAddress: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48",
		DstTokenAddress: "0x2791bca1f2de4661ed88a30c99a7a9449aa84174",
		Amount:          "100000000",
		WalletAddress:   "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb",
		EnableEstimate:  true,
	}

	customPreset := quoter.CustomPreset{
		AuctionDuration:    300,
		AuctionStartAmount: "100000000",
		AuctionEndAmount:   "99500000",
	}

	quote, err := sdk.GetQuoteWithCustomPreset(context.Background(), params, customPreset)
	require.NoError(t, err)
	assert.NotNil(t, quote)
}

func TestSDK_CreateOrder(t *testing.T) {
	config := Config{
		URL:     "https://api.example.com",
		AuthKey: "test-key",
	}
	sdk := NewSDK(config)

	quote := sampleQuote()
	secret := repeatHex("00", 31) + "01"
	hashLock, err := hashlock.ForSingleFill(secret)
	require.NoError(t, err)

	secretHashes := []string{repeatHex("00", 32)}
	params := OrderParams{
		WalletAddress: "0x0742D35CC6634c0532925A3b844bc9E7595f0Beb",
		HashLock:      hashLock,
		SecretHashes:  secretHashes,
		Preset:        "fast",
	}

	preparedOrder, err := sdk.CreateOrder(context.Background(), quote, params)
	require.NoError(t, err)
	assert.NotNil(t, preparedOrder)
	assert.NotEmpty(t, preparedOrder.Hash)
	assert.Equal(t, quote.QuoteID, preparedOrder.QuoteID)
}

func TestSDK_CreateOrder_NoQuoteID(t *testing.T) {
	config := Config{
		URL:     "https://api.example.com",
		AuthKey: "test-key",
	}
	sdk := NewSDK(config)

	quote := sampleQuote()
	quote.QuoteID = "" // No quote ID

	secret := repeatHex("00", 31) + "01"
	hashLock, err := hashlock.ForSingleFill(secret)
	require.NoError(t, err)

	params := OrderParams{
		WalletAddress: "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb",
		HashLock:      hashLock,
		SecretHashes:  []string{},
		Preset:        "fast",
	}

	preparedOrder, err := sdk.CreateOrder(context.Background(), quote, params)
	assert.Error(t, err)
	assert.Nil(t, preparedOrder)
}

func TestSDK_SignOrder(t *testing.T) {
	mockProvider := &testutils.MockBlockchainProvider{}
	mockProvider.On("SignTypedData", context.Background(), mock.Anything, mock.Anything).Return("0xsignature", nil)

	config := Config{
		URL:                "https://api.example.com",
		AuthKey:            "test-key",
		BlockchainProvider: mockProvider,
	}
	sdk := NewSDK(config)

	order := sampleEvmOrder()
	signature, err := sdk.SignOrder(context.Background(), order, chains.Ethereum)
	require.NoError(t, err)
	assert.NotEmpty(t, signature)
	mockProvider.AssertExpectations(t)
}

func TestSDK_SignOrder_NoProvider(t *testing.T) {
	config := Config{
		URL:     "https://api.example.com",
		AuthKey: "test-key",
	}
	sdk := NewSDK(config)

	order := sampleEvmOrder()
	signature, err := sdk.SignOrder(context.Background(), order, chains.Ethereum)
	assert.Error(t, err)
	assert.Empty(t, signature)
}

func TestSDK_CreateOrder_ExtensionBuilt(t *testing.T) {
	config := Config{
		URL:     "https://api.example.com",
		AuthKey: "test-key",
	}
	sdk := NewSDK(config)

	quote := sampleQuote()
	secret := repeatHex("00", 31) + "01"
	hashLock, err := hashlock.ForSingleFill(secret)
	require.NoError(t, err)

	params := OrderParams{
		WalletAddress: "0x0742D35CC6634c0532925A3b844bc9E7595f0Beb",
		HashLock:      hashLock,
		SecretHashes:  []string{repeatHex("00", 32)},
		Preset:        "fast",
	}

	preparedOrder, err := sdk.CreateOrder(context.Background(), quote, params)
	require.NoError(t, err)
	assert.NotNil(t, preparedOrder)

	order, ok := preparedOrder.Order.(*orders.EvmCrossChainOrder)
	require.True(t, ok, "order should be EvmCrossChainOrder")
	assert.NotEmpty(t, order.Extension, "extension should be built")
	assert.Equal(t, "0x", order.Extension[:2], "extension should start with 0x")
}

func TestSDK_CreateOrder_ExtensionWithWhitelist(t *testing.T) {
	config := Config{
		URL:     "https://api.example.com",
		AuthKey: "test-key",
	}
	sdk := NewSDK(config)

	quote := sampleQuote()
	quote.Whitelist = []string{
		"0x1111111111111111111111111111111111111111",
		"0x2222222222222222222222222222222222222222",
	}

	secret := repeatHex("00", 31) + "01"
	hashLock, err := hashlock.ForSingleFill(secret)
	require.NoError(t, err)

	params := OrderParams{
		WalletAddress: "0x0742D35CC6634c0532925A3b844bc9E7595f0Beb",
		HashLock:      hashLock,
		SecretHashes:  []string{repeatHex("00", 32)},
		Preset:        "fast",
	}

	preparedOrder, err := sdk.CreateOrder(context.Background(), quote, params)
	require.NoError(t, err)
	assert.NotNil(t, preparedOrder)

	order, ok := preparedOrder.Order.(*orders.EvmCrossChainOrder)
	require.True(t, ok)
	assert.NotEmpty(t, order.Extension)
	assert.NotEmpty(t, order.Whitelist)
	assert.Equal(t, len(quote.Whitelist), len(order.Whitelist))
}

func TestSDK_SignNativeOrder(t *testing.T) {
	config := Config{
		URL:     "https://api.example.com",
		AuthKey: "test-key",
	}
	sdk := NewSDK(config)

	order := sampleEvmOrder()
	maker := addresses.NewEvmAddress("0x0742D35CC6634c0532925A3b844bc9E7595f0Beb")
	signature := sdk.SignNativeOrder(order, maker)
	assert.NotEmpty(t, signature)
	assert.Equal(t, 132, len(signature)) // 0x + 65 bytes * 2 hex chars
}

func TestSDK_SubmitOrder(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
	}))
	defer server.Close()

	mockProvider := &testutils.MockBlockchainProvider{}
	mockProvider.On("SignTypedData", context.Background(), mock.AnythingOfType("string"), mock.AnythingOfType("eip712.TypedData")).Return("0xsignature", nil)

	config := Config{
		URL:                server.URL,
		AuthKey:            "test-key",
		BlockchainProvider: mockProvider,
	}
	sdk := NewSDK(config)

	order := sampleEvmOrder()
	secretHashes := []string{repeatHex("00", 32)}

	orderInfo, err := sdk.SubmitOrder(context.Background(), chains.Ethereum, order, "quote-123", secretHashes)
	require.NoError(t, err)
	assert.NotNil(t, orderInfo)
	assert.NotEmpty(t, orderInfo.OrderHash)
	mockProvider.AssertExpectations(t)
}

func TestSDK_SubmitOrder_InvalidSecretHashes(t *testing.T) {
	config := Config{
		URL:     "https://api.example.com",
		AuthKey: "test-key",
	}
	sdk := NewSDK(config)

	order := sampleEvmOrder()
	order.MultipleFillsAllowed = false
	secretHashes := []string{"0x1", "0x2"} // Multiple when not allowed

	orderInfo, err := sdk.SubmitOrder(context.Background(), chains.Ethereum, order, "quote-123", secretHashes)
	assert.Error(t, err)
	assert.Nil(t, orderInfo)
}

func TestSDK_AnnounceOrder(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
	}))
	defer server.Close()

	config := Config{
		URL:     server.URL,
		AuthKey: "test-key",
	}
	sdk := NewSDK(config)

	orderHash := make([]byte, 32)
	secret := repeatHex("00", 31) + "01"
	hashLock, err := hashlock.ForSingleFill(secret)
	require.NoError(t, err)

	auctionDetails := &auction.AuctionDetails{
		StartTime:       big.NewInt(1735689600),
		Duration:        big.NewInt(300),
		InitialRateBump: 100,
		Points:          []auction.AuctionPoint{},
		GasCost: auction.GasCost{
			GasBumpEstimate:  big.NewInt(100000),
			GasPriceEstimate: big.NewInt(2000000000),
		},
	}

	order := &orders.SolanaCrossChainOrder{
		OrderHash:            orderHash,
		HashLock:             hashLock,
		Auction:              auctionDetails,
		MultipleFillsAllowed: true,
	}

	secretHashes := []string{repeatHex("00", 32)}
	hash, err := sdk.AnnounceOrder(context.Background(), order, "quote-123", secretHashes)
	require.NoError(t, err)
	assert.NotEmpty(t, hash)
}

func TestSDK_PlaceOrder(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
	}))
	defer server.Close()

	mockProvider := &testutils.MockBlockchainProvider{}
	mockProvider.On("SignTypedData", context.Background(), mock.AnythingOfType("string"), mock.AnythingOfType("eip712.TypedData")).Return("0xsignature", nil)

	config := Config{
		URL:                server.URL,
		AuthKey:            "test-key",
		BlockchainProvider: mockProvider,
	}
	sdk := NewSDK(config)

	quote := sampleQuote()
	secret := repeatHex("00", 31) + "01"
	hashLock, err := hashlock.ForSingleFill(secret)
	require.NoError(t, err)

	secretHashes := []string{repeatHex("00", 32)}
	params := OrderParams{
		WalletAddress: "0x0742D35CC6634c0532925A3b844bc9E7595f0Beb",
		HashLock:      hashLock,
		SecretHashes:  secretHashes,
		Preset:        "fast",
	}

	orderInfo, err := sdk.PlaceOrder(context.Background(), quote, params)
	require.NoError(t, err)
	assert.NotNil(t, orderInfo)
	mockProvider.AssertExpectations(t)
}

func TestSDK_PlaceOrder_SolanaOrder(t *testing.T) {
	config := Config{
		URL:     "https://api.example.com",
		AuthKey: "test-key",
	}
	sdk := NewSDK(config)

	quote := sampleQuote()
	quote.SrcChainID = chains.Solana // Solana source

	secret := repeatHex("00", 31) + "01"
	hashLock, err := hashlock.ForSingleFill(secret)
	require.NoError(t, err)

	params := OrderParams{
		WalletAddress: "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGxrHJf99",
		HashLock:      hashLock,
		SecretHashes:  []string{},
		Preset:        "fast",
		Receiver:      "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb",
	}

	orderInfo, err := sdk.PlaceOrder(context.Background(), quote, params)
	assert.Error(t, err) // Should error because Solana orders need AnnounceOrder
	assert.Nil(t, orderInfo)
}

func TestSDK_GetOrderStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		_ = json.NewEncoder(w).Encode(apiorders.OrderStatusResponse{
			OrderHash:  "0x123",
			Status:     apiorders.OrderStatusPending,
			SrcChainID: chains.Ethereum,
			DstChainID: chains.Polygon,
		})
	}))
	defer server.Close()

	config := Config{
		URL:     server.URL,
		AuthKey: "test-key",
	}
	sdk := NewSDK(config)

	status, err := sdk.GetOrderStatus(context.Background(), "0x123")
	require.NoError(t, err)
	assert.NotNil(t, status)
	assert.Equal(t, "0x123", status.OrderHash)
}

func TestSDK_GetActiveOrders(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		_ = json.NewEncoder(w).Encode(apiorders.ActiveOrdersResponse{
			Items:      []apiorders.ActiveOrder{},
			TotalCount: 0,
		})
	}))
	defer server.Close()

	config := Config{
		URL:     server.URL,
		AuthKey: "test-key",
	}
	sdk := NewSDK(config)

	params := apiorders.ActiveOrdersRequestParams{
		Page:  1,
		Limit: 10,
	}

	orders, err := sdk.GetActiveOrders(context.Background(), params)
	require.NoError(t, err)
	assert.NotNil(t, orders)
}

func TestSDK_SubmitSecret(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
	}))
	defer server.Close()

	config := Config{
		URL:     server.URL,
		AuthKey: "test-key",
	}
	sdk := NewSDK(config)

	err := sdk.SubmitSecret(context.Background(), "0x123", "0xsecret")
	assert.NoError(t, err)
}

func TestSDK_BuildCancelOrderCallData(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		_ = json.NewEncoder(w).Encode(apiorders.OrderStatusResponse{
			OrderHash:  "0x123",
			Status:     apiorders.OrderStatusPending,
			SrcChainID: chains.Ethereum,
			DstChainID: chains.Polygon,
			Order: &orders.LimitOrderV4Struct{
				MakerTraits: "12345",
			},
		})
	}))
	defer server.Close()

	config := Config{
		URL:     server.URL,
		AuthKey: "test-key",
	}
	sdk := NewSDK(config)

	callData, err := sdk.BuildCancelOrderCallData(context.Background(), "0x123")
	require.NoError(t, err)
	assert.NotEmpty(t, callData)
	assert.True(t, len(callData) > 2)
	assert.Equal(t, "0x", callData[:2])
}

func TestSDK_BuildCancelOrderCallData_NonEvm(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		_ = json.NewEncoder(w).Encode(apiorders.OrderStatusResponse{
			OrderHash:  "0x123",
			Status:     apiorders.OrderStatusPending,
			SrcChainID: chains.Solana, // Non-EVM
			DstChainID: chains.Ethereum,
		})
	}))
	defer server.Close()

	config := Config{
		URL:     server.URL,
		AuthKey: "test-key",
	}
	sdk := NewSDK(config)

	callData, err := sdk.BuildCancelOrderCallData(context.Background(), "0x123")
	assert.Error(t, err)
	assert.Empty(t, callData)
}

func TestSDK_GenerateSecrets(t *testing.T) {
	secrets, err := GenerateSecrets(1)
	require.NoError(t, err)
	assert.Equal(t, 1, len(secrets))
	assert.True(t, len(secrets[0]) > 2)
	assert.Equal(t, "0x", secrets[0][:2])

	secrets2, err := GenerateSecrets(4)
	require.NoError(t, err)
	assert.Equal(t, 4, len(secrets2))
	for _, secret := range secrets2 {
		assert.True(t, len(secret) > 2)
		assert.Equal(t, "0x", secret[:2])
	}
}

func TestSDK_GenerateSecrets_Zero(t *testing.T) {
	secrets, err := GenerateSecrets(0)
	require.NoError(t, err)
	assert.Equal(t, 0, len(secrets))
}
