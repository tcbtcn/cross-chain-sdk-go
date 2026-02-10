package orders

import (
	"math/big"
	"strings"
	"testing"

	"github.com/dawitel/cross-chain-sdk-go/chains"
	"github.com/dawitel/cross-chain-sdk-go/domains/addresses"
	"github.com/dawitel/cross-chain-sdk-go/domains/hashlock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func sampleEvmOrder() *EvmCrossChainOrder {
	maker := addresses.NewEvmAddress("0x0742D35CC6634c0532925A3b844bc9E7595f0Beb")
	makerAsset := addresses.NewEvmAddress("0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48")
	takerAsset := addresses.NewEvmAddress("0x2791bca1f2de4661ed88a30c99a7a9449aa84174")
	receiver := maker

	secret := "0x0000000000000000000000000000000000000000000000000000000000000001"
	hashLock, _ := hashlock.ForSingleFill(secret)

	return &EvmCrossChainOrder{
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

func TestEvmCrossChainOrder_Build(t *testing.T) {
	order := sampleEvmOrder()

	built, err := order.Build()
	require.NoError(t, err)
	assert.NotNil(t, built)
	assert.Equal(t, order.Maker.ToString(), built.Maker)
	assert.Equal(t, order.MakerAsset.ToString(), built.MakerAsset)
	assert.Equal(t, order.MakingAmount.String(), built.MakingAmount)
	assert.Equal(t, order.TakingAmount.String(), built.TakingAmount)
}

func TestEvmCrossChainOrder_GetOrderHash(t *testing.T) {
	order := sampleEvmOrder()

	hash, err := order.GetOrderHash(chains.Ethereum)
	require.NoError(t, err)
	require.NotEmpty(t, hash, "GetOrderHash should return a non-empty hash")
	require.True(t, len(hash) > 2, "Hash should have at least 2 characters for 0x prefix")
	assert.Equal(t, "0x", hash[:2])

	// Should be deterministic
	hash2, err := order.GetOrderHash(chains.Ethereum)
	require.NoError(t, err)
	assert.Equal(t, hash, hash2)
}

func TestEvmCrossChainOrder_GetTypedData(t *testing.T) {
	order := sampleEvmOrder()

	typedData, err := order.GetTypedData(chains.Ethereum)
	require.NoError(t, err)
	assert.NotNil(t, typedData)

	typedDataMap, ok := typedData.(map[string]interface{})
	require.True(t, ok)
	assert.NotNil(t, typedDataMap["types"])
	assert.NotNil(t, typedDataMap["domain"])
	assert.NotNil(t, typedDataMap["message"])
	assert.Equal(t, "Order", typedDataMap["primaryType"])
}

func TestEvmCrossChainOrder_GetTypedData_DifferentChains(t *testing.T) {
	order := sampleEvmOrder()

	typedData1, err1 := order.GetTypedData(chains.Ethereum)
	require.NoError(t, err1)
	typedData2, err2 := order.GetTypedData(chains.Polygon)
	require.NoError(t, err2)

	// Domain chainId should be different
	td1, _ := typedData1.(map[string]interface{})
	td2, _ := typedData2.(map[string]interface{})

	domain1 := td1["domain"].(map[string]interface{})
	domain2 := td2["domain"].(map[string]interface{})

	assert.NotEqual(t, domain1["chainId"], domain2["chainId"])
}

func TestSolanaCrossChainOrder_GetOrderHash(t *testing.T) {
	orderHash := make([]byte, 32)
	for i := range orderHash {
		orderHash[i] = byte(i)
	}

	secret := "0x0000000000000000000000000000000000000000000000000000000000000001"
	hashLock, err := hashlock.ForSingleFill(secret)
	require.NoError(t, err)

	order := &SolanaCrossChainOrder{
		OrderHash:            orderHash,
		HashLock:             hashLock,
		MultipleFillsAllowed: true,
	}

	hash := order.GetOrderHash()
	assert.NotEmpty(t, hash)
}

func TestSolanaCrossChainOrder_ToJSON(t *testing.T) {
	orderHash := make([]byte, 32)
	secret := "0x0000000000000000000000000000000000000000000000000000000000000001"
	hashLock, err := hashlock.ForSingleFill(secret)
	require.NoError(t, err)

	// Create a minimal valid Solana order with required fields
	srcToken, err := addresses.NewSolanaAddress("So11111111111111111111111111111111111111112")
	require.NoError(t, err)
	dstToken, err := addresses.EvmAddressFromString("0x0000000000000000000000000000000000000000")
	require.NoError(t, err)
	maker, err := addresses.NewSolanaAddress("So11111111111111111111111111111111111111112")
	require.NoError(t, err)
	receiver, err := addresses.EvmAddressFromString("0x0000000000000000000000000000000000000000")
	require.NoError(t, err)

	order := &SolanaCrossChainOrder{
		OrderHash:            orderHash,
		HashLock:             hashLock,
		MultipleFillsAllowed: true,
		SrcToken:             srcToken,
		DstToken:             dstToken,
		Maker:                maker,
		Receiver:             receiver,
		SrcAmount:            big.NewInt(1000000),
		MinDstAmount:         big.NewInt(950000),
		SrcSafetyDeposit:     big.NewInt(1000),
		DstSafetyDeposit:     big.NewInt(1000),
		TimeLocks:            nil, // Can be nil for test
		DstChainID:           chains.Ethereum,
		Salt:                 big.NewInt(12345),
		Source:               "sdk",
		SrcAssetIsNative:     false,
	}

	jsonData := order.ToJSON()
	assert.NotNil(t, jsonData)
	
	// Check that all required top-level keys exist
	_, hasOrderInfo := jsonData["orderInfo"]
	assert.True(t, hasOrderInfo, "json should contain orderInfo")
	
	_, hasEscrowParams := jsonData["escrowParams"]
	assert.True(t, hasEscrowParams, "json should contain escrowParams")
	
	_, hasDetails := jsonData["details"]
	assert.True(t, hasDetails, "json should contain details")
	
	_, hasExtra := jsonData["extra"]
	assert.True(t, hasExtra, "json should contain extra")
	
	orderInfo, ok := jsonData["orderInfo"].(map[string]interface{})
	assert.True(t, ok, "orderInfo should be a map")
	if ok {
		_, hasSrcToken := orderInfo["srcToken"]
		assert.True(t, hasSrcToken, "orderInfo should contain srcToken")
		
		_, hasDstToken := orderInfo["dstToken"]
		assert.True(t, hasDstToken, "orderInfo should contain dstToken")
	}
}

func TestLimitOrderV4Struct(t *testing.T) {
	struct_ := LimitOrderV4Struct{
		Maker:         "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb",
		MakerAsset:    "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48",
		TakerAsset:    "0x2791bca1f2de4661ed88a30c99a7a9449aa84174",
		MakingAmount:  "100000000",
		TakingAmount:  "99500000",
		Receiver:      "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb",
		AllowedSender: "0x0000000000000000000000000000000000000000",
		MakerTraits:   "0",
		Salt:          "12345",
		Expiration:    "1735689600",
		Nonce:         "67890",
	}

	assert.NotEmpty(t, struct_.Maker)
	assert.NotEmpty(t, struct_.MakerAsset)
}

func TestEvmCrossChainOrder_WithNativeAsset(t *testing.T) {
	maker := addresses.NewEvmAddress("0x0742D35CC6634c0532925A3b844bc9E7595f0Beb")
	makerAsset := addresses.EvmNative
	takerAsset := addresses.NewEvmAddress("0x2791bca1f2de4661ed88a30c99a7a9449aa84174")
	receiver := maker

	secret := "0x0000000000000000000000000000000000000000000000000000000000000001"
	hashLock, err := hashlock.ForSingleFill(secret)
	require.NoError(t, err)

	order := &EvmCrossChainOrder{
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

	built, err := order.Build()
	require.NoError(t, err)
	assert.True(t, order.MakerAsset.IsNative())
	// Build() returns the address as returned by ToString(), which may have mixed case
	// Compare case-insensitively since addresses are case-insensitive
	assert.Equal(t, strings.ToLower(addresses.NativeCurrency), strings.ToLower(built.MakerAsset))
}

func TestEvmCrossChainOrder_Build_NilMaker(t *testing.T) {
	order := sampleEvmOrder()
	order.Maker = nil

	_, err := order.Build()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "maker address is required")
}

func TestEvmCrossChainOrder_Build_NilMakerAsset(t *testing.T) {
	order := sampleEvmOrder()
	order.MakerAsset = nil

	_, err := order.Build()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "makerAsset address is required")
}

func TestEvmCrossChainOrder_Build_NilTakerAsset(t *testing.T) {
	order := sampleEvmOrder()
	order.TakerAsset = nil

	_, err := order.Build()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "takerAsset address is required")
}

func TestEvmCrossChainOrder_Build_NilReceiver(t *testing.T) {
	order := sampleEvmOrder()
	order.Receiver = nil

	_, err := order.Build()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "receiver address is required")
}

func TestEvmCrossChainOrder_Build_NilHashLock(t *testing.T) {
	order := sampleEvmOrder()
	order.HashLock = nil

	_, err := order.Build()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "hashLock is required")
}

func TestEvmCrossChainOrder_GetTypedData_NilMaker(t *testing.T) {
	order := sampleEvmOrder()
	order.Maker = nil

	_, err := order.GetTypedData(chains.Ethereum)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "maker address is required")
}

func TestEvmCrossChainOrder_GetOrderHash_NilMaker(t *testing.T) {
	order := sampleEvmOrder()
	order.Maker = nil

	_, err := order.GetOrderHash(chains.Ethereum)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "maker address is required")
}
