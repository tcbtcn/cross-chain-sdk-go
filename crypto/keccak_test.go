package crypto

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestKeccak256(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		want  string
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
			want:  "99ff0d9125e1fc9531a11262e15aeb2c60509a078c4cc4c64cefdfb06ff68647",
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
	input := []byte("hello")
	expected := "0x1c8aff950685c2ed4bc3174f3472287b56d9517b9c948127319a09a7a36deac8"
	got := Keccak256Hex(input)
	assert.Equal(t, expected, got)
}

func TestKeccak256FromHex(t *testing.T) {
	tests := []struct {
		name    string
		hexStr  string
		wantErr bool
	}{
		{
			name:    "valid hex",
			hexStr:  "0x68656c6c6f",
			wantErr: false,
		},
		{
			name:    "valid hex without 0x",
			hexStr:  "68656c6c6f",
			wantErr: true, // Our function requires 0x prefix
		},
		{
			name:    "invalid hex",
			hexStr:  "0xzzzz",
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
			values:  []interface{}{uint64(123), "0x0000000000000000000000000000000000000000000000000000000000000001"},
			wantErr: false,
		},
		{
			name:    "multiple uint64",
			types:   []string{"uint64", "uint64", "uint64"},
			values:  []interface{}{uint64(1), uint64(2), uint64(3)},
			wantErr: false,
		},
		{
			name:    "mismatched types and values",
			types:   []string{"uint64", "uint64"},
			values:  []interface{}{uint64(1)},
			wantErr: true,
		},
		{
			name:    "unsupported type",
			types:   []string{"string"},
			values:  []interface{}{"hello"},
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
			}
		})
	}
}

func TestSolidityPackedKeccak256Hex(t *testing.T) {
	types := []string{"uint64", "uint64"}
	values := []interface{}{uint64(1), uint64(2)}
	hash, err := SolidityPackedKeccak256(types, values)
	require.NoError(t, err)
	expected := "0x" + hex.EncodeToString(hash)
	got, err := SolidityPackedKeccak256Hex(types, values)
	require.NoError(t, err)
	assert.Equal(t, expected, got)
}

func TestKeccak256_Deterministic(t *testing.T) {
	input := []byte("test input")
	hash1 := Keccak256(input)
	hash2 := Keccak256(input)
	assert.Equal(t, hash1, hash2)
}
