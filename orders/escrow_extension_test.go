package orders

import (
	"math/big"
	"testing"

	"github.com/dawitel/cross-chain-sdk-go/chains"
	"github.com/dawitel/cross-chain-sdk-go/domains/addresses"
	"github.com/dawitel/cross-chain-sdk-go/domains/auction"
	"github.com/dawitel/cross-chain-sdk-go/domains/hashlock"
	"github.com/dawitel/cross-chain-sdk-go/domains/timelocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestEscrowExtension_Build(t *testing.T) {
	escrowFactory := addresses.NewEvmAddress("0x1111111111111111111111111111111111111111")

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

	whitelist := []AuctionWhitelistItem{
		{
			Address:   addresses.NewEvmAddress("0x2222222222222222222222222222222222222222"),
			AllowFrom: big.NewInt(1000),
		},
	}

	postInteractionData := NewSettlementPostInteractionData(whitelist, big.NewInt(1735689600))

	hashLock, err := hashlock.ForSingleFill("0x0000000000000000000000000000000000000000000000000000000000000001")
	require.NoError(t, err)

	dstToken := addresses.NewEvmAddress("0x2791bca1f2de4661ed88a30c99a7a9449aa84174")
	srcSafetyDeposit := big.NewInt(1000000)
	dstSafetyDeposit := big.NewInt(2000000)

	tl, err := timelocks.NewTimeLocks(timelocks.TimeLocksParams{
		SrcWithdrawal:         big.NewInt(3600),
		SrcPublicWithdrawal:   big.NewInt(7200),
		SrcCancellation:       big.NewInt(10800),
		SrcPublicCancellation: big.NewInt(14400),
		DstWithdrawal:         big.NewInt(3600),
		DstPublicWithdrawal:   big.NewInt(7200),
		DstCancellation:       big.NewInt(10800),
	})
	require.NoError(t, err)

	dstAddressFirstPart := addresses.ZeroComplement

	escrowExt := NewEscrowExtension(
		escrowFactory,
		auctionDetails,
		postInteractionData,
		nil,
		hashLock,
		chains.Polygon,
		dstToken,
		srcSafetyDeposit,
		dstSafetyDeposit,
		tl,
		dstAddressFirstPart,
	)

	extension, err := escrowExt.Build()
	require.NoError(t, err)
	assert.NotEmpty(t, extension)
	assert.Equal(t, "0x", extension[:2])
	assert.True(t, len(extension) > 2)
}

func TestEscrowExtension_Build_NilAddress(t *testing.T) {
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

	postInteractionData := NewSettlementPostInteractionData([]AuctionWhitelistItem{}, big.NewInt(1735689600))
	hashLock, _ := hashlock.ForSingleFill("0x0000000000000000000000000000000000000000000000000000000000000001")
	dstToken := addresses.NewEvmAddress("0x2791bca1f2de4661ed88a30c99a7a9449aa84174")
	tl, _ := timelocks.NewTimeLocks(timelocks.TimeLocksParams{
		SrcWithdrawal:         big.NewInt(3600),
		SrcPublicWithdrawal:   big.NewInt(7200),
		SrcCancellation:       big.NewInt(10800),
		SrcPublicCancellation: big.NewInt(14400),
		DstWithdrawal:         big.NewInt(3600),
		DstPublicWithdrawal:   big.NewInt(7200),
		DstCancellation:       big.NewInt(10800),
	})

	escrowExt := NewEscrowExtension(
		nil,
		auctionDetails,
		postInteractionData,
		nil,
		hashLock,
		chains.Polygon,
		dstToken,
		big.NewInt(1000000),
		big.NewInt(2000000),
		tl,
		addresses.ZeroComplement,
	)

	_, err := escrowExt.Build()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "address is required")
}

func TestEscrowExtension_Build_WithCustomData(t *testing.T) {
	escrowFactory := addresses.NewEvmAddress("0x1111111111111111111111111111111111111111")

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

	postInteractionData := NewSettlementPostInteractionData([]AuctionWhitelistItem{}, big.NewInt(1735689600))
	hashLock, _ := hashlock.ForSingleFill("0x0000000000000000000000000000000000000000000000000000000000000001")
	dstToken := addresses.NewEvmAddress("0x2791bca1f2de4661ed88a30c99a7a9449aa84174")
	tl, _ := timelocks.NewTimeLocks(timelocks.TimeLocksParams{
		SrcWithdrawal:         big.NewInt(3600),
		SrcPublicWithdrawal:   big.NewInt(7200),
		SrcCancellation:       big.NewInt(10800),
		SrcPublicCancellation: big.NewInt(14400),
		DstWithdrawal:         big.NewInt(3600),
		DstPublicWithdrawal:   big.NewInt(7200),
		DstCancellation:       big.NewInt(10800),
	})

	customData := addresses.NewAddressComplement(big.NewInt(12345))

	escrowExt := NewEscrowExtension(
		escrowFactory,
		auctionDetails,
		postInteractionData,
		nil,
		hashLock,
		chains.Polygon,
		dstToken,
		big.NewInt(1000000),
		big.NewInt(2000000),
		tl,
		customData,
	)

	extension, err := escrowExt.Build()
	require.NoError(t, err)
	assert.NotEmpty(t, extension)
}

func TestEscrowExtension_Build_Deterministic(t *testing.T) {
	escrowFactory := addresses.NewEvmAddress("0x1111111111111111111111111111111111111111")

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

	postInteractionData := NewSettlementPostInteractionData([]AuctionWhitelistItem{}, big.NewInt(1735689600))
	hashLock, _ := hashlock.ForSingleFill("0x0000000000000000000000000000000000000000000000000000000000000001")
	dstToken := addresses.NewEvmAddress("0x2791bca1f2de4661ed88a30c99a7a9449aa84174")
	tl, _ := timelocks.NewTimeLocks(timelocks.TimeLocksParams{
		SrcWithdrawal:         big.NewInt(3600),
		SrcPublicWithdrawal:   big.NewInt(7200),
		SrcCancellation:       big.NewInt(10800),
		SrcPublicCancellation: big.NewInt(14400),
		DstWithdrawal:         big.NewInt(3600),
		DstPublicWithdrawal:   big.NewInt(7200),
		DstCancellation:       big.NewInt(10800),
	})

	escrowExt := NewEscrowExtension(
		escrowFactory,
		auctionDetails,
		postInteractionData,
		nil,
		hashLock,
		chains.Polygon,
		dstToken,
		big.NewInt(1000000),
		big.NewInt(2000000),
		tl,
		addresses.ZeroComplement,
	)

	extension1, err := escrowExt.Build()
	require.NoError(t, err)

	extension2, err := escrowExt.Build()
	require.NoError(t, err)

	assert.Equal(t, extension1, extension2)
}

func TestEscrowExtension_Build_NativeToken(t *testing.T) {
	escrowFactory := addresses.NewEvmAddress("0x1111111111111111111111111111111111111111")

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

	postInteractionData := NewSettlementPostInteractionData([]AuctionWhitelistItem{}, big.NewInt(1735689600))
	hashLock, _ := hashlock.ForSingleFill("0x0000000000000000000000000000000000000000000000000000000000000001")
	dstToken := addresses.EvmNative
	tl, _ := timelocks.NewTimeLocks(timelocks.TimeLocksParams{
		SrcWithdrawal:         big.NewInt(3600),
		SrcPublicWithdrawal:   big.NewInt(7200),
		SrcCancellation:       big.NewInt(10800),
		SrcPublicCancellation: big.NewInt(14400),
		DstWithdrawal:         big.NewInt(3600),
		DstPublicWithdrawal:   big.NewInt(7200),
		DstCancellation:       big.NewInt(10800),
	})

	escrowExt := NewEscrowExtension(
		escrowFactory,
		auctionDetails,
		postInteractionData,
		nil,
		hashLock,
		chains.Polygon,
		dstToken,
		big.NewInt(1000000),
		big.NewInt(2000000),
		tl,
		addresses.ZeroComplement,
	)

	extension, err := escrowExt.Build()
	require.NoError(t, err)
	assert.NotEmpty(t, extension)
}
