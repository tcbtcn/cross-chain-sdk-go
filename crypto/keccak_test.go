package crypto

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test vectors from known Keccak256 hashes
func TestKeccak256(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		want  string // hex encoded expected hash
	}{
		{
			name:  "empty input",
			input: []byte{},
			want:  "c5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470",
		},
		{
			name:  "simple string",
			input: []byte("hello"),
			want:  "1c8aff950685c2ed4bc3174f3472287b56d9517b9c948127319a09a7a36deac8",
		},
		{
			name:  "zero bytes",
			input: []byte{0x00, 0x00, 0x00},
			want:  "99ff0d9125e1fc9531a11262e15aeb2c60509a078c4cc4c6c4c4efdfb06ff68647",
		},
		{
			name:  "32 bytes of zeros",
			input: make([]byte, 32),
			want:  "290decd9548b62a8d60345a988386fc84ba6bc95484008f6362f93160ef3e563",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Keccak256(tt.input)
			wantBytes, err := hex.DecodeString(tt.want)
			require.NoError(t, err)
			assert.Equal(t, wantBytes, got)
		})
	}
}

func TestKeccak256Hex(t *testing.T) {
	input := []byte("test")
	got := Keccak256Hex(input)

	expectedHash := Keccak256(input)
	expectedHex := "0x" + hex.EncodeToString(expectedHash)

	assert.Equal(t, expectedHex, got)
	assert.True(t, len(got) > 2)
	assert.Equal(t, "0x", got[:2])
}

func TestKeccak256FromHex(t *testing.T) {
	tests := []struct {
		name    string
		hexStr  string
		wantErr bool
	}{
		{
			name:    "valid hex",
			hexStr:  "0x68656c6c6f", // "hello"
			wantErr: false,
		},
		{
			name:    "valid hex without 0x",
			hexStr:  "68656c6c6f",
			wantErr: true, // function expects 0x prefix
		},
		{
			name:    "invalid hex",
			hexStr:  "0xgh",
			wantErr: true,
		},
		{
			name:    "empty hex",
			hexStr:  "0x",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Keccak256FromHex(tt.hexStr)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.Equal(t, 32, len(got)) // Keccak256 always returns 32 bytes
			}
		})
	}
}

func TestSolidityPackedKeccak256(t *testing.T) {
	tests := []struct {
		name    string
		types   []string
		values  []interface{}
		wantErr bool
	}{
		{
			name:    "uint64 and bytes32",
			types:   []string{"uint64", "bytes32"},
			values:  []interface{}{uint64(123), repeatHex("00", 32)},
			wantErr: false,
		},
		{
			name:    "multiple uint64",
			types:   []string{"uint64", "uint64"},
			values:  []interface{}{uint64(1), uint64(2)},
			wantErr: false,
		},
		{
			name:    "mismatched types and values",
			types:   []string{"uint64"},
			values:  []interface{}{uint64(1), uint64(2)},
			wantErr: false, // function doesn't check length mismatch
		},
		{
			name:    "unsupported type",
			types:   []string{"string"},
			values:  []interface{}{"test"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SolidityPackedKeccak256(tt.types, tt.values)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, got)
				assert.Equal(t, 32, len(got))
			}
		})
	}
}

func TestSolidityPackedKeccak256Hex(t *testing.T) {
	types := []string{"uint64"}
	values := []interface{}{uint64(123)}

	got, err := SolidityPackedKeccak256Hex(types, values)
	require.NoError(t, err)

	assert.True(t, len(got) > 2)
	assert.Equal(t, "0x", got[:2])

	// Verify it matches the bytes version
	bytesHash, err := SolidityPackedKeccak256(types, values)
	require.NoError(t, err)
	expectedHex := "0x" + hex.EncodeToString(bytesHash)
	assert.Equal(t, expectedHex, got)
}

func TestKeccak256_Deterministic(t *testing.T) {
	input := []byte("deterministic test")
	hash1 := Keccak256(input)
	hash2 := Keccak256(input)
	assert.Equal(t, hash1, hash2)
}
