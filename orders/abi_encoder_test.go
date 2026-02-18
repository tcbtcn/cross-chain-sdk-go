package orders

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tcbtcn/cross-chain-sdk-go/chains"
	"github.com/tcbtcn/cross-chain-sdk-go/domains/addresses"
	"github.com/tcbtcn/cross-chain-sdk-go/domains/hashlock"
	"github.com/tcbtcn/cross-chain-sdk-go/domains/timelocks"
)

func TestEncodeExtraEscrowData(t *testing.T) {
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

	encoded, err := EncodeExtraEscrowData(
		hashLock,
		int(chains.Polygon),
		dstToken,
		srcSafetyDeposit,
		dstSafetyDeposit,
		tl,
	)
	require.NoError(t, err)
	assert.NotEmpty(t, encoded)
	assert.NotEqual(t, "0x", encoded)
}

func TestEncodeExtraEscrowData_NilHashLock(t *testing.T) {
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

	_, err = EncodeExtraEscrowData(
		nil,
		int(chains.Polygon),
		dstToken,
		srcSafetyDeposit,
		dstSafetyDeposit,
		tl,
	)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "hashLock is required")
}

func TestEncodeExtraEscrowData_NilDstToken(t *testing.T) {
	hashLock, err := hashlock.ForSingleFill("0x0000000000000000000000000000000000000000000000000000000000000001")
	require.NoError(t, err)

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

	_, err = EncodeExtraEscrowData(
		hashLock,
		int(chains.Polygon),
		nil,
		srcSafetyDeposit,
		dstSafetyDeposit,
		tl,
	)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "dstToken is required")
}

func TestEncodeExtraEscrowData_ExceedsUint128Max(t *testing.T) {
	hashLock, err := hashlock.ForSingleFill("0x0000000000000000000000000000000000000000000000000000000000000001")
	require.NoError(t, err)

	dstToken := addresses.NewEvmAddress("0x2791bca1f2de4661ed88a30c99a7a9449aa84174")
	uint128Max, _ := new(big.Int).SetString("0xffffffffffffffffffffffffffffffff", 0)
	srcSafetyDeposit := new(big.Int).Add(uint128Max, big.NewInt(1))
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

	_, err = EncodeExtraEscrowData(
		hashLock,
		int(chains.Polygon),
		dstToken,
		srcSafetyDeposit,
		dstSafetyDeposit,
		tl,
	)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "srcSafetyDeposit exceeds UINT_128_MAX")
}

func TestEncodeExtraEscrowData_NativeToken(t *testing.T) {
	hashLock, err := hashlock.ForSingleFill("0x0000000000000000000000000000000000000000000000000000000000000001")
	require.NoError(t, err)

	dstToken := addresses.EvmNative
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

	encoded, err := EncodeExtraEscrowData(
		hashLock,
		int(chains.Polygon),
		dstToken,
		srcSafetyDeposit,
		dstSafetyDeposit,
		tl,
	)
	require.NoError(t, err)
	assert.NotEmpty(t, encoded)
}

func TestEncodeExtraEscrowData_Deterministic(t *testing.T) {
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

	encoded1, err := EncodeExtraEscrowData(
		hashLock,
		int(chains.Polygon),
		dstToken,
		srcSafetyDeposit,
		dstSafetyDeposit,
		tl,
	)
	require.NoError(t, err)

	encoded2, err := EncodeExtraEscrowData(
		hashLock,
		int(chains.Polygon),
		dstToken,
		srcSafetyDeposit,
		dstSafetyDeposit,
		tl,
	)
	require.NoError(t, err)

	assert.Equal(t, encoded1, encoded2)
}
