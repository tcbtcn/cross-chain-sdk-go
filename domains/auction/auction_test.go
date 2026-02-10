package auction

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAuctionDetails_HashForSolana(t *testing.T) {
	ad := &AuctionDetails{
		StartTime:       big.NewInt(1000),
		InitialRateBump: 100,
		Duration:        big.NewInt(300),
		Points: []AuctionPoint{
			{Delay: 0, ToTokenAmount: "1000000"},
			{Delay: 100, ToTokenAmount: "950000"},
		},
		GasCost: GasCost{
			GasBumpEstimate:  big.NewInt(100000),
			GasPriceEstimate: big.NewInt(20000000000),
		},
	}

	hash, err := ad.HashForSolana()
	require.NoError(t, err)
	assert.Equal(t, 32, len(hash))

	// Should be deterministic
	hash2, err := ad.HashForSolana()
	require.NoError(t, err)
	assert.Equal(t, hash, hash2)
}

func TestAuctionDetails_HashForSolanaHex(t *testing.T) {
	ad := &AuctionDetails{
		StartTime:       big.NewInt(1000),
		InitialRateBump: 100,
		Duration:        big.NewInt(300),
		Points:          []AuctionPoint{},
		GasCost: GasCost{
			GasBumpEstimate:  big.NewInt(100000),
			GasPriceEstimate: big.NewInt(20000000000),
		},
	}

	hashHex, err := ad.HashForSolanaHex()
	require.NoError(t, err)
	assert.True(t, len(hashHex) > 2)
	assert.Equal(t, "0x", hashHex[:2])
}

func TestAuctionDetails_Encode(t *testing.T) {
	ad := &AuctionDetails{
		StartTime:       big.NewInt(1000),
		InitialRateBump: 100,
		Duration:        big.NewInt(300),
		Points: []AuctionPoint{
			{Delay: 0, ToTokenAmount: "1000000"},
		},
		GasCost: GasCost{
			GasBumpEstimate:  big.NewInt(100000),
			GasPriceEstimate: big.NewInt(20000000000),
		},
	}

	encoded, err := ad.encode()
	require.NoError(t, err)
	assert.NotNil(t, encoded)
}

func TestAuctionDetails_Encode_InvalidAmount(t *testing.T) {
	ad := &AuctionDetails{
		StartTime:       big.NewInt(1000),
		InitialRateBump: 100,
		Duration:        big.NewInt(300),
		Points: []AuctionPoint{
			{Delay: 0, ToTokenAmount: "invalid"},
		},
		GasCost: GasCost{
			GasBumpEstimate:  big.NewInt(100000),
			GasPriceEstimate: big.NewInt(20000000000),
		},
	}

	_, err := ad.encode()
	assert.Error(t, err)
}

func TestAuctionDetails_EmptyPoints(t *testing.T) {
	ad := &AuctionDetails{
		StartTime:       big.NewInt(1000),
		InitialRateBump: 100,
		Duration:        big.NewInt(300),
		Points:          []AuctionPoint{},
		GasCost: GasCost{
			GasBumpEstimate:  big.NewInt(100000),
			GasPriceEstimate: big.NewInt(20000000000),
		},
	}

	hash, err := ad.HashForSolana()
	require.NoError(t, err)
	assert.Equal(t, 32, len(hash))
}

func TestAuctionDetails_MultiplePoints(t *testing.T) {
	ad := &AuctionDetails{
		StartTime:       big.NewInt(1000),
		InitialRateBump: 100,
		Duration:        big.NewInt(300),
		Points: []AuctionPoint{
			{Delay: 0, ToTokenAmount: "1000000"},
			{Delay: 50, ToTokenAmount: "975000"},
			{Delay: 100, ToTokenAmount: "950000"},
			{Delay: 150, ToTokenAmount: "925000"},
		},
		GasCost: GasCost{
			GasBumpEstimate:  big.NewInt(100000),
			GasPriceEstimate: big.NewInt(20000000000),
		},
	}

	hash, err := ad.HashForSolana()
	require.NoError(t, err)
	assert.Equal(t, 32, len(hash))
}
