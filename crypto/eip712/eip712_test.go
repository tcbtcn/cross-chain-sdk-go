package eip712

import (
	"math/big"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func repeatHex(s string, n int) string {
	return "0x" + strings.Repeat(s, n)
}

func TestHashTypedData(t *testing.T) {
	typedData := TypedData{
		Types: map[string]interface{}{
			"EIP712Domain": []interface{}{
				map[string]interface{}{"name": "name", "type": "string"},
				map[string]interface{}{"name": "version", "type": "string"},
				map[string]interface{}{"name": "chainId", "type": "uint256"},
			},
			"Message": []interface{}{
				map[string]interface{}{"name": "value", "type": "uint256"},
			},
		},
		PrimaryType: "Message",
		Domain: map[string]interface{}{
			"name":    "Test",
			"version": "1",
			"chainId": "1",
		},
		Message: map[string]interface{}{
			"value": "100",
		},
	}

	hash, err := HashTypedData(typedData)
	require.NoError(t, err)
	assert.Equal(t, 32, len(hash))
}

func TestTypedData_HashStruct(t *testing.T) {
	td := &TypedData{
		Types: map[string]interface{}{
			"TestStruct": []interface{}{
				map[string]interface{}{"name": "field1", "type": "uint256"},
				map[string]interface{}{"name": "field2", "type": "address"},
			},
		},
	}

	data := map[string]interface{}{
		"field1": "100",
		"field2": "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb",
	}

	hash, err := td.HashStruct("TestStruct", data)
	require.NoError(t, err)
	assert.Equal(t, 32, len(hash))
}

func TestTypedData_EncodeField(t *testing.T) {
	td := &TypedData{
		Types: map[string]interface{}{},
	}

	tests := []struct {
		name      string
		fieldType string
		value     interface{}
		wantErr   bool
	}{
		{
			name:      "string type",
			fieldType: "string",
			value:     "test",
			wantErr:   false,
		},
		{
			name:      "bytes type",
			fieldType: "bytes",
			value:     []byte{0x12, 0x34},
			wantErr:   false,
		},
		{
			name:      "bytes32 type from string",
			fieldType: "bytes32",
			value:     repeatHex("00", 32),
			wantErr:   false,
		},
		{
			name:      "bytes32 type from bytes",
			fieldType: "bytes32",
			value:     make([]byte, 32),
			wantErr:   false,
		},
		{
			name:      "address type from string",
			fieldType: "address",
			value:     "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb",
			wantErr:   false,
		},
		{
			name:      "address type from Address",
			fieldType: "address",
			value:     common.HexToAddress("0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb"),
			wantErr:   false,
		},
		{
			name:      "uint256 type from string",
			fieldType: "uint256",
			value:     "100",
			wantErr:   false,
		},
		{
			name:      "uint256 type from big.Int",
			fieldType: "uint256",
			value:     big.NewInt(100),
			wantErr:   false,
		},
		{
			name:      "uint256 type from int64",
			fieldType: "uint256",
			value:     int64(100),
			wantErr:   false,
		},
		{
			name:      "uint64 type from uint64",
			fieldType: "uint64",
			value:     uint64(100),
			wantErr:   false,
		},
		{
			name:      "invalid string type",
			fieldType: "string",
			value:     123,
			wantErr:   true,
		},
		{
			name:      "unsupported type",
			fieldType: "unknown",
			value:     "test",
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			encoded, err := td.EncodeField(tt.fieldType, tt.value)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, encoded)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, encoded)
			}
		})
	}
}

func TestFormatSignature(t *testing.T) {
	tests := []struct {
		name string
		sig  []byte
		want string
	}{
		{
			name: "valid signature",
			sig:  make([]byte, 65),
			want: repeatHex("00", 65),
		},
		{
			name: "signature with recovery >= 27",
			sig:  append(make([]byte, 64), byte(27)),
			want: repeatHex("00", 64) + "00",
		},
		{
			name: "invalid length",
			sig:  []byte{0x00},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FormatSignature(tt.sig)
			if tt.want == "" {
				assert.Empty(t, got)
			} else {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
