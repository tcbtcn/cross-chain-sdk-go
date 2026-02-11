package orders

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	"github.com/dawitel/cross-chain-sdk-go/domains/addresses"
	"github.com/dawitel/cross-chain-sdk-go/domains/hashlock"
	"github.com/dawitel/cross-chain-sdk-go/domains/timelocks"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common/math"
)

const (
	UINT_128_MAX = "0xffffffffffffffffffffffffffffffff"
)

var (
	extraDataTypes = []string{
		"bytes32", // hashLock
		"uint256", // dstChainId
		"uint256", // dstToken
		"uint256", // safety deposits (packed)
		"uint256", // timeLocks
	}
)

func EncodeExtraEscrowData(
	hashLock *hashlock.HashLock,
	dstChainID int,
	dstToken addresses.AddressLike,
	srcSafetyDeposit *big.Int,
	dstSafetyDeposit *big.Int,
	timeLocks *timelocks.TimeLocks,
) (string, error) {
	if hashLock == nil {
		return "", fmt.Errorf("hashLock is required")
	}
	if dstToken == nil {
		return "", fmt.Errorf("dstToken is required")
	}
	if srcSafetyDeposit == nil {
		return "", fmt.Errorf("srcSafetyDeposit is required")
	}
	if dstSafetyDeposit == nil {
		return "", fmt.Errorf("dstSafetyDeposit is required")
	}
	if timeLocks == nil {
		return "", fmt.Errorf("timeLocks is required")
	}

	hashLockBytes, err := hashLock.ToBuffer()
	if err != nil {
		return "", fmt.Errorf("failed to get hashLock buffer: %w", err)
	}

	dstTokenAddr, ok := dstToken.(*addresses.EvmAddress)
	if !ok {
		return "", fmt.Errorf("dstToken must be EvmAddress")
	}

	dstTokenNativeAsZero := dstTokenAddr.NativeAsZero().(*addresses.EvmAddress)
	dstTokenBigInt := dstTokenNativeAsZero.ToBigInt()

	uint128Max, _ := new(big.Int).SetString(UINT_128_MAX[2:], 16)
	if srcSafetyDeposit.Cmp(uint128Max) > 0 {
		return "", fmt.Errorf("srcSafetyDeposit exceeds UINT_128_MAX")
	}
	if dstSafetyDeposit.Cmp(uint128Max) > 0 {
		return "", fmt.Errorf("dstSafetyDeposit exceeds UINT_128_MAX")
	}

	packedSafetyDeposits := new(big.Int).Lsh(srcSafetyDeposit, 128)
	packedSafetyDeposits.Or(packedSafetyDeposits, dstSafetyDeposit)

	timeLocksBigInt := timeLocks.Build()

	arguments := abi.Arguments{}
	for _, t := range extraDataTypes {
		ty, err := abi.NewType(t, "", nil)
		if err != nil {
			return "", fmt.Errorf("failed to create type %s: %w", t, err)
		}
		arguments = append(arguments, abi.Argument{Type: ty})
	}

	values := []interface{}{
		[32]byte(hashLockBytes),
		big.NewInt(int64(dstChainID)),
		dstTokenBigInt,
		packedSafetyDeposits,
		timeLocksBigInt,
	}

	encoded, err := arguments.Pack(values...)
	if err != nil {
		return "", fmt.Errorf("failed to pack arguments: %w", err)
	}

	result := "0x" + hex.EncodeToString(encoded)
	return strings.TrimPrefix(result, "0x"), nil
}

func EncodeUint256(value *big.Int) []byte {
	return math.U256Bytes(value)
}

func EncodeAddress(addr *addresses.EvmAddress) []byte {
	return addr.ToBuffer()
}

func EncodeBytes32(data []byte) []byte {
	if len(data) > 32 {
		panic("data exceeds 32 bytes")
	}
	result := make([]byte, 32)
	copy(result[32-len(data):], data)
	return result
}
