package addresses

import (
	"math/big"
	"testing"

	"github.com/gagliardetto/solana-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSolanaAddress(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{
			name:    "valid address",
			value:   "So11111111111111111111111111111111111111112",
			wantErr: false,
		},
		{
			name:    "system program",
			value:   "11111111111111111111111111111111",
			wantErr: false,
		},
		{
			name:    "invalid address",
			value:   "invalid",
			wantErr: true,
		},
		{
			name:    "too short",
			value:   "short",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addr, err := NewSolanaAddress(tt.value)
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

func TestSolanaAddressFromString(t *testing.T) {
	addr, err := SolanaAddressFromString("So11111111111111111111111111111111111111112")
	require.NoError(t, err)
	assert.NotNil(t, addr)
}

func TestSolanaAddressFromBuffer(t *testing.T) {
	tests := []struct {
		name    string
		buf     []byte
		wantErr bool
	}{
		{
			name:    "valid 32 bytes",
			buf:     make([]byte, 32),
			wantErr: false,
		},
		{
			name:    "too short",
			buf:     make([]byte, 31),
			wantErr: true,
		},
		{
			name:    "too long",
			buf:     make([]byte, 33),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			addr, err := SolanaAddressFromBuffer(tt.buf)
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

func TestSolanaAddressFromPublicKey(t *testing.T) {
	pubkey := solana.PublicKey{}
	addr := SolanaAddressFromPublicKey(pubkey)
	assert.NotNil(t, addr)
}

func TestSolanaAddressFromBigInt(t *testing.T) {
	val := big.NewInt(123456789)
	addr := SolanaAddressFromBigInt(val)
	assert.NotNil(t, addr)
}

func TestSolanaAddress_ToString(t *testing.T) {
	addr, err := NewSolanaAddress("So11111111111111111111111111111111111111112")
	require.NoError(t, err)
	str := addr.ToString()
	assert.NotEmpty(t, str)
}

func TestSolanaAddress_ToBuffer(t *testing.T) {
	addr, err := NewSolanaAddress("So11111111111111111111111111111111111111112")
	require.NoError(t, err)
	buf := addr.ToBuffer()
	assert.Equal(t, 32, len(buf))
}

func TestSolanaAddress_Equal(t *testing.T) {
	addr1, err1 := NewSolanaAddress("So11111111111111111111111111111111111111112")
	require.NoError(t, err1)
	addr2, err2 := NewSolanaAddress("So11111111111111111111111111111111111111112")
	require.NoError(t, err2)
	addr3, err3 := NewSolanaAddress("11111111111111111111111111111111")
	require.NoError(t, err3)

	assert.True(t, addr1.Equal(addr2))
	assert.False(t, addr1.Equal(addr3))
}

func TestSolanaAddress_IsZero(t *testing.T) {
	zeroAddr := SolanaAddressFromBigInt(big.NewInt(0))
	assert.True(t, zeroAddr.IsZero())

	nonZeroAddr, err := NewSolanaAddress("So11111111111111111111111111111111111111112")
	require.NoError(t, err)
	assert.False(t, nonZeroAddr.IsZero())
}

func TestSolanaConstants(t *testing.T) {
	assert.NotNil(t, SolanaAssociatedTokenProgramID)
	assert.NotNil(t, SolanaTokenProgramID)
	assert.NotNil(t, SolanaSystemProgramID)
	assert.NotNil(t, SolanaWrappedNative)
	assert.NotNil(t, SolanaZero)
	assert.True(t, SolanaZero.IsZero())
}
