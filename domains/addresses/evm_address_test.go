package addresses

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewEvmAddress(t *testing.T) {
	addr := NewEvmAddress("0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb")
	assert.NotNil(t, addr)
	// Address is normalized to checksummed format by go-ethereum
	assert.Equal(t, "0x0742D35CC6634c0532925A3b844bc9E7595f0Beb", addr.ToString())
}

func TestEvmAddressFromString(t *testing.T) {
	tests := []struct {
		name    string
		address string
		wantErr bool
	}{
		{
			name:    "valid address",
			address: "0x0742D35CC6634c0532925A3b844bc9E7595f0Beb",
			wantErr: false,
		},
		{
			name:    "zero address",
			address: ZeroAddress,
			wantErr: false,
		},
		{
			name:    "native currency",
			address: NativeCurrency,
			wantErr: false,
		},
		{
			name:    "invalid address",
			address: "invalid",
			wantErr: true,
		},
		{
			name:    "too short",
			address: "0x123",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addr, err := EvmAddressFromString(tt.address)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, addr)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, addr)
			}
		})
	}
}

func TestEvmAddressFromBigInt(t *testing.T) {
	val := big.NewInt(123456789)
	addr := EvmAddressFromBigInt(val)
	assert.NotNil(t, addr)
}

func TestEvmAddressFromBuffer(t *testing.T) {
	tests := []struct {
		name    string
		buf     []byte
		wantErr bool
	}{
		{
			name:    "valid 20 bytes",
			buf:     make([]byte, 20),
			wantErr: false,
		},
		{
			name:    "valid more than 20 bytes",
			buf:     make([]byte, 32),
			wantErr: false,
		},
		{
			name:    "too short",
			buf:     make([]byte, 19),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addr, err := EvmAddressFromBuffer(tt.buf)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, addr)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, addr)
			}
		})
	}
}

func TestEvmAddress_IsZero(t *testing.T) {
	zeroAddr := NewEvmAddress(ZeroAddress)
	assert.True(t, zeroAddr.IsZero())

	nonZeroAddr := NewEvmAddress("0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb")
	assert.False(t, nonZeroAddr.IsZero())
}

func TestEvmAddress_IsNative(t *testing.T) {
	nativeAddr := NewEvmAddress(NativeCurrency)
	assert.True(t, nativeAddr.IsNative())

	nonNativeAddr := NewEvmAddress("0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb")
	assert.False(t, nonNativeAddr.IsNative())
}

func TestEvmAddress_NativeAsZero(t *testing.T) {
	nativeAddr := NewEvmAddress(NativeCurrency)
	result := nativeAddr.NativeAsZero()
	assert.True(t, result.IsZero())

	regularAddr := NewEvmAddress("0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb")
	result = regularAddr.NativeAsZero()
	assert.Equal(t, regularAddr.ToString(), result.ToString())
}

func TestEvmAddress_ZeroAsNative(t *testing.T) {
	zeroAddr := NewEvmAddress(ZeroAddress)
	result := zeroAddr.ZeroAsNative()
	assert.True(t, result.IsNative())

	regularAddr := NewEvmAddress("0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb")
	result = regularAddr.ZeroAsNative()
	assert.Equal(t, regularAddr.ToString(), result.ToString())
}

func TestEvmAddress_Equal(t *testing.T) {
	addr1 := NewEvmAddress("0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb")
	addr2 := NewEvmAddress("0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb")
	addr3 := NewEvmAddress("0x8ba1f109551bD432803012645Hac136c22C9")

	assert.True(t, addr1.Equal(addr2))
	assert.False(t, addr1.Equal(addr3))
}

func TestEvmAddress_ToString(t *testing.T) {
	addr := NewEvmAddress("0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb")
	// go-ethereum normalizes addresses to checksummed format
	assert.Equal(t, "0x0742D35CC6634c0532925A3b844bc9E7595f0Beb", addr.ToString())
}

func TestEvmAddress_ToBuffer(t *testing.T) {
	addr := NewEvmAddress("0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb")
	buf := addr.ToBuffer()
	assert.Equal(t, 20, len(buf))
}

func TestEvmAddress_ToBigInt(t *testing.T) {
	addr := NewEvmAddress("0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb")
	val := addr.ToBigInt()
	assert.NotNil(t, val)
}

func TestEvmZeroAndNative(t *testing.T) {
	assert.True(t, EvmZero.IsZero())
	assert.True(t, EvmNative.IsNative())
}
