package crypto

import (
	"math/big"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func repeatHex(s string, n int) string {
	return "0x" + strings.Repeat(s, n)
}

func TestEncodeCancelOrder(t *testing.T) {
	tests := []struct {
		name        string
		orderHash   string
		makerTraits *big.Int
		wantErr     bool
	}{
		{
			name:        "valid order hash and traits",
			orderHash:   repeatHex("00", 32),
			makerTraits: big.NewInt(12345),
			wantErr:     false,
		},
		{
			name:        "zero maker traits",
			orderHash:   repeatHex("00", 32),
			makerTraits: big.NewInt(0),
			wantErr:     false,
		},
		{
			name:        "large maker traits",
			orderHash:   repeatHex("ff", 32),
			makerTraits: mustParseBigInt(t, "115792089237316195423570985008687907853269984665640564039457584007913129639935"),
			wantErr:     false,
		},
		{
			name:        "valid order hash format",
			orderHash:   "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
			makerTraits: big.NewInt(100),
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := EncodeCancelOrder(tt.orderHash, tt.makerTraits)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, got)
				assert.True(t, len(got) > 2)
				assert.Equal(t, "0x", got[:2])
				// Method ID (4 bytes) + makerTraits (32 bytes) + orderHash (32 bytes) = 68 bytes = 136 hex chars + "0x"
				assert.Equal(t, 138, len(got))
			}
		})
	}
}

func TestEncodeCancelOrder_Format(t *testing.T) {
	orderHash := repeatHex("00", 32)
	makerTraits := big.NewInt(12345)

	callData, err := EncodeCancelOrder(orderHash, makerTraits)
	require.NoError(t, err)

	// Should start with 0x
	assert.True(t, len(callData) > 2)
	assert.Equal(t, "0x", callData[:2])

	// Should have method ID (4 bytes) + makerTraits (32 bytes) + orderHash (32 bytes)
	hexData := callData[2:]
	assert.Equal(t, 136, len(hexData)) // 68 bytes * 2 hex chars
}

func mustParseBigInt(t *testing.T, s string) *big.Int {
	val, ok := new(big.Int).SetString(s, 10)
	require.True(t, ok)
	return val
}
