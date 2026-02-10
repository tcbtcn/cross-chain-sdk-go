package utils

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsHexBytes(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{
			name:  "valid hex with 0x",
			input: "0x1234",
			want:  true,
		},
		{
			name:  "valid hex lowercase",
			input: "0xabcdef",
			want:  true,
		},
		{
			name:  "valid hex uppercase",
			input: "0xABCDEF",
			want:  true,
		},
		{
			name:  "valid hex mixed case",
			input: "0xAbCdEf",
			want:  true,
		},
		{
			name:  "valid empty hex",
			input: "0x",
			want:  true,
		},
		{
			name:  "invalid without 0x",
			input: "1234",
			want:  false,
		},
		{
			name:  "invalid empty string",
			input: "",
			want:  false,
		},
		{
			name:  "invalid single character",
			input: "0",
			want:  false,
		},
		{
			name:  "invalid odd length hex",
			input: "0x123",
			want:  false,
		},
		{
			name:  "invalid non-hex characters",
			input: "0xghij",
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsHexBytes(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestBufferFromHex(t *testing.T) {
	tests := []struct {
		name      string
		hexStr    string
		bytesSize int
		want      []byte
		wantErr   bool
	}{
		{
			name:      "valid hex with 0x",
			hexStr:    "0x1234",
			bytesSize: -1,
			want:      []byte{0x12, 0x34},
			wantErr:   false,
		},
		{
			name:      "invalid hex without 0x",
			hexStr:    "1234",
			bytesSize: -1,
			want:      nil,
			wantErr:   true,
		},
		{
			name:      "valid hex with padding",
			hexStr:    "0x12",
			bytesSize: 4,
			want:      []byte{0x00, 0x00, 0x00, 0x12},
			wantErr:   false,
		},
		{
			name:      "valid hex exact size",
			hexStr:    "0x1234",
			bytesSize: 2,
			want:      []byte{0x12, 0x34},
			wantErr:   false,
		},
		{
			name:      "invalid hex exceeds size",
			hexStr:    "0x123456",
			bytesSize: 2,
			want:      nil,
			wantErr:   true,
		},
		{
			name:      "invalid hex string",
			hexStr:    "0xgh",
			bytesSize: -1,
			want:      nil,
			wantErr:   true,
		},
		{
			name:      "invalid without 0x",
			hexStr:    "1234",
			bytesSize: -1,
			want:      nil,
			wantErr:   true,
		},
		{
			name:      "empty hex",
			hexStr:    "0x",
			bytesSize: -1,
			want:      []byte{},
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := BufferFromHex(tt.hexStr, tt.bytesSize)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestBufferToHex(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		want  string
	}{
		{
			name:  "empty buffer",
			input: []byte{},
			want:  "0x",
		},
		{
			name:  "single byte",
			input: []byte{0x12},
			want:  "0x12",
		},
		{
			name:  "multiple bytes",
			input: []byte{0x12, 0x34, 0x56},
			want:  "0x123456",
		},
		{
			name:  "zero bytes",
			input: []byte{0x00, 0x00},
			want:  "0x0000",
		},
		{
			name:  "large buffer",
			input: []byte{0xff, 0xee, 0xdd, 0xcc, 0xbb, 0xaa},
			want:  "0xffeeddccbbaa",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BufferToHex(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGetBytesCount(t *testing.T) {
	tests := []struct {
		name    string
		hexStr  string
		want    int64
		wantErr bool
	}{
		{
			name:    "valid 2 bytes",
			hexStr:  "0x1234",
			want:    2,
			wantErr: false,
		},
		{
			name:    "valid 4 bytes",
			hexStr:  "0x12345678",
			want:    4,
			wantErr: false,
		},
		{
			name:    "valid 32 bytes",
			hexStr:  repeatHex("00", 32),
			want:    32,
			wantErr: false,
		},
		{
			name:    "valid empty",
			hexStr:  "0x",
			want:    0,
			wantErr: false,
		},
		{
			name:    "invalid without 0x",
			hexStr:  "1234",
			want:    0,
			wantErr: true,
		},
		{
			name:    "invalid non-hex",
			hexStr:  "0xgh",
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetBytesCount(tt.hexStr)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Equal(t, int64(0), got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestBufferFromHexRoundTrip(t *testing.T) {
	original := []byte{0x12, 0x34, 0x56, 0x78, 0x9a, 0xbc, 0xde, 0xf0}
	hexStr := BufferToHex(original)
	decoded, err := BufferFromHex(hexStr, -1)
	require.NoError(t, err)
	assert.Equal(t, original, decoded)
}

func repeatHex(s string, n int) string {
	return "0x" + strings.Repeat(s, n)
}
