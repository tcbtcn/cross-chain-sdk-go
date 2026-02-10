package testutils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
	"testing"

	"github.com/dawitel/cross-chain-sdk-go/chains"
	"github.com/dawitel/cross-chain-sdk-go/domains/addresses"
	"github.com/dawitel/cross-chain-sdk-go/domains/hashlock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// GenerateRandomSecret generates a random 32-byte secret
func GenerateRandomSecret(t *testing.T) string {
	secret := make([]byte, 32)
	_, err := rand.Read(secret)
	require.NoError(t, err)
	return "0x" + hex.EncodeToString(secret)
}

// GenerateRandomSecrets generates n random secrets
func GenerateRandomSecrets(t *testing.T, n int) []string {
	secrets := make([]string, n)
	for i := 0; i < n; i++ {
		secrets[i] = GenerateRandomSecret(t)
	}
	return secrets
}

// CreateTestHashLock creates a hash lock for testing
func CreateTestHashLock(t *testing.T, secrets []string) *hashlock.HashLock {
	if len(secrets) == 1 {
		hl, err := hashlock.ForSingleFill(secrets[0])
		require.NoError(t, err)
		return hl
	}

	leaves, err := hashlock.GetMerkleLeaves(secrets)
	require.NoError(t, err)
	hl, err := hashlock.ForMultipleFills(leaves)
	require.NoError(t, err)
	return hl
}

// CreateTestEvmAddress creates a test EVM address from string
func CreateTestEvmAddress(t *testing.T, addr string) *addresses.EvmAddress {
	evmAddr, err := addresses.EvmAddressFromString(addr)
	require.NoError(t, err)
	return evmAddr
}

// AssertBigIntEqual asserts two big.Ints are equal
func AssertBigIntEqual(t *testing.T, expected, actual *big.Int, msgAndArgs ...interface{}) {
	assert.Equal(t, expected.String(), actual.String(), msgAndArgs...)
}

// AssertHashLockEqual asserts two HashLocks are equal
func AssertHashLockEqual(t *testing.T, expected, actual *hashlock.HashLock) {
	assert.Equal(t, expected.ToString(), actual.ToString())
}

// CreateTestQuoteParams creates test quote parameters
func CreateTestQuoteParams(srcChain, dstChain chains.SupportedChain) map[string]interface{} {
	return map[string]interface{}{
		"srcChainId":      srcChain,
		"dstChainId":      dstChain,
		"srcTokenAddress": "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48",
		"dstTokenAddress": "0x2791bca1f2de4661ed88a30c99a7a9449aa84174",
		"amount":          "100000000",
		"walletAddress":   "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb",
		"enableEstimate":  true,
	}
}

// MustParseBigInt parses a string to big.Int or fails the test
func MustParseBigInt(t *testing.T, s string) *big.Int {
	val, ok := new(big.Int).SetString(s, 10)
	require.True(t, ok, fmt.Sprintf("failed to parse big.Int from %s", s))
	return val
}

// MustParseHex parses a hex string to bytes or fails the test
func MustParseHex(t *testing.T, hexStr string) []byte {
	if len(hexStr) >= 2 && hexStr[0:2] == "0x" {
		hexStr = hexStr[2:]
	}
	data, err := hex.DecodeString(hexStr)
	require.NoError(t, err, fmt.Sprintf("failed to parse hex: %s", hexStr))
	return data
}
