package sdk

import (
	"math/big"
	"testing"

	"github.com/dawitel/cross-chain-sdk-go/chains"
	"github.com/dawitel/cross-chain-sdk-go/domains/hashlock"
	"github.com/dawitel/cross-chain-sdk-go/orders"
	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	config := Config{
		URL:     "https://api.example.com",
		AuthKey: "test-key",
	}

	assert.Equal(t, "https://api.example.com", config.URL)
	assert.Equal(t, "test-key", config.AuthKey)
}

func TestQuoteParams(t *testing.T) {
	params := QuoteParams{
		SrcChainID:      chains.Ethereum,
		DstChainID:      chains.Polygon,
		SrcTokenAddress: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48",
		DstTokenAddress: "0x2791bca1f2de4661ed88a30c99a7a9449aa84174",
		Amount:          "100000000",
		WalletAddress:   "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb",
		EnableEstimate:  true,
		Permit:          "0x123",
		TakingFeeBps:    100,
		Source:          "test",
		IsPermit2:       true,
	}

	assert.Equal(t, chains.Ethereum, params.SrcChainID)
	assert.Equal(t, chains.Polygon, params.DstChainID)
	assert.True(t, params.EnableEstimate)
}

func TestOrderParams(t *testing.T) {
	hashLock, _ := hashlock.ForSingleFill(repeatHex("00", 32))

	params := OrderParams{
		WalletAddress: "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb",
		HashLock:      hashLock,
		SecretHashes:  []string{repeatHex("00", 32)},
		Preset:        "fast",
		Nonce:         big.NewInt(12345),
		Receiver:      "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb",
	}

	assert.NotNil(t, params.HashLock)
	assert.Equal(t, "fast", params.Preset)
	assert.NotNil(t, params.Nonce)
}

func TestOrderInfo(t *testing.T) {
	orderInfo := OrderInfo{
		Order:     orders.LimitOrderV4Struct{},
		Signature: "0xsignature",
		QuoteID:   "quote-123",
		OrderHash: "0xhash",
		Extension: "0xextension",
	}

	assert.Equal(t, "0xsignature", orderInfo.Signature)
	assert.Equal(t, "quote-123", orderInfo.QuoteID)
	assert.Equal(t, "0xhash", orderInfo.OrderHash)
}

func TestPreparedOrder(t *testing.T) {
	preparedOrder := PreparedOrder{
		Order:   &orders.EvmCrossChainOrder{},
		Hash:    "0xhash",
		QuoteID: "quote-123",
	}

	assert.NotNil(t, preparedOrder.Order)
	assert.Equal(t, "0xhash", preparedOrder.Hash)
	assert.Equal(t, "quote-123", preparedOrder.QuoteID)
}

func TestTakingFeeInfo(t *testing.T) {
	fee := TakingFeeInfo{
		TakingFeeBps:      100,
		TakingFeeReceiver: "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb",
	}

	assert.Equal(t, 100, fee.TakingFeeBps)
	assert.NotEmpty(t, fee.TakingFeeReceiver)
}
