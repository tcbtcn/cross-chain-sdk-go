package auction

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuctionDetails_EncodeForEvm(t *testing.T) {
	ad := &AuctionDetails{
		StartTime:       big.NewInt(1735689600),
		InitialRateBump: 100,
		Duration:        big.NewInt(300),
		Points: []AuctionPoint{
			{Delay: 0, ToTokenAmount: "10000"},
			{Delay: 50, ToTokenAmount: "5000"},
		},
		GasCost: GasCost{
			GasBumpEstimate:  big.NewInt(100000),
			GasPriceEstimate: big.NewInt(2000000000),
		},
	}

	encoded, err := ad.EncodeForEvm()
	require.NoError(t, err)
	assert.NotNil(t, encoded)
	assert.NotEmpty(t, encoded)
}

func TestAuctionDetails_EncodeForEvm_EmptyPoints(t *testing.T) {
	ad := &AuctionDetails{
		StartTime:       big.NewInt(1735689600),
		InitialRateBump: 100,
		Duration:        big.NewInt(300),
		Points:          []AuctionPoint{},
		GasCost: GasCost{
			GasBumpEstimate:  big.NewInt(100000),
			GasPriceEstimate: big.NewInt(2000000000),
		},
	}

	encoded, err := ad.EncodeForEvm()
	require.NoError(t, err)
	assert.NotNil(t, encoded)
}

func TestAuctionDetails_EncodeForEvm_ExceedsUint24(t *testing.T) {
	ad := &AuctionDetails{
		StartTime:       big.NewInt(1735689600),
		InitialRateBump: 0x1000000,
		Duration:        big.NewInt(300),
		Points:          []AuctionPoint{},
		GasCost: GasCost{
			GasBumpEstimate:  big.NewInt(100000),
			GasPriceEstimate: big.NewInt(2000000000),
		},
	}

	_, err := ad.EncodeForEvm()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "initialRateBump exceeds uint24")
}

func TestAuctionDetails_EncodeForEvm_ExceedsUint32(t *testing.T) {
	ad := &AuctionDetails{
		StartTime:       big.NewInt(0x100000000),
		InitialRateBump: 100,
		Duration:        big.NewInt(300),
		Points:          []AuctionPoint{},
		GasCost: GasCost{
			GasBumpEstimate:  big.NewInt(100000),
			GasPriceEstimate: big.NewInt(2000000000),
		},
	}

	_, err := ad.EncodeForEvm()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "startTime exceeds uint32")
}

func TestAuctionDetails_EncodeForEvm_InvalidToTokenAmount(t *testing.T) {
	ad := &AuctionDetails{
		StartTime:       big.NewInt(1735689600),
		InitialRateBump: 100,
		Duration:        big.NewInt(300),
		Points: []AuctionPoint{
			{Delay: 0, ToTokenAmount: "invalid"},
		},
		GasCost: GasCost{
			GasBumpEstimate:  big.NewInt(100000),
			GasPriceEstimate: big.NewInt(2000000000),
		},
	}

	_, err := ad.EncodeForEvm()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid toTokenAmount")
}

func TestAuctionDetails_EncodeForEvm_CoefficientExceedsUint16(t *testing.T) {
	ad := &AuctionDetails{
		StartTime:       big.NewInt(1735689600),
		InitialRateBump: 100,
		Duration:        big.NewInt(300),
		Points: []AuctionPoint{
			{Delay: 0, ToTokenAmount: "70000"},
		},
		GasCost: GasCost{
			GasBumpEstimate:  big.NewInt(100000),
			GasPriceEstimate: big.NewInt(2000000000),
		},
	}

	_, err := ad.EncodeForEvm()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "coefficient exceeds uint16")
}

func TestAuctionDetails_EncodeForEvm_Deterministic(t *testing.T) {
	ad := &AuctionDetails{
		StartTime:       big.NewInt(1735689600),
		InitialRateBump: 100,
		Duration:        big.NewInt(300),
		Points: []AuctionPoint{
			{Delay: 0, ToTokenAmount: "10000"},
		},
		GasCost: GasCost{
			GasBumpEstimate:  big.NewInt(100000),
			GasPriceEstimate: big.NewInt(2000000000),
		},
	}

	encoded1, err := ad.EncodeForEvm()
	require.NoError(t, err)

	encoded2, err := ad.EncodeForEvm()
	require.NoError(t, err)

	assert.Equal(t, encoded1, encoded2)
}

func TestAuctionDetails_EncodeForEvm_MultiplePoints(t *testing.T) {
	ad := &AuctionDetails{
		StartTime:       big.NewInt(1735689600),
		InitialRateBump: 100,
		Duration:        big.NewInt(300),
		Points: []AuctionPoint{
			{Delay: 0, ToTokenAmount: "10000"},
			{Delay: 50, ToTokenAmount: "5000"},
			{Delay: 100, ToTokenAmount: "2500"},
		},
		GasCost: GasCost{
			GasBumpEstimate:  big.NewInt(100000),
			GasPriceEstimate: big.NewInt(2000000000),
		},
	}

	encoded, err := ad.EncodeForEvm()
	require.NoError(t, err)
	assert.NotNil(t, encoded)
}
