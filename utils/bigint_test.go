package utils

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsBigIntString(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  bool
	}{
		{
			name:  "valid positive number",
			input: "123456789",
			want:  true,
		},
		{
			name:  "valid zero",
			input: "0",
			want:  true,
		},
		{
			name:  "valid large number",
			input: "999999999999999999999999999999999999999999",
			want:  true,
		},
		{
			name:  "invalid empty string",
			input: "",
			want:  false,
		},
		{
			name:  "invalid non-numeric",
			input: "abc",
			want:  false,
		},
		{
			name:  "invalid hex string",
			input: "0x123",
			want:  false,
		},
		{
			name:  "invalid with decimal",
			input: "123.456",
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsBigIntString(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestBigIntFromString(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *big.Int
		wantErr bool
	}{
		{
			name:    "valid positive number",
			input:   "123456789",
			want:    big.NewInt(123456789),
			wantErr: false,
		},
		{
			name:    "valid zero",
			input:   "0",
			want:    big.NewInt(0),
			wantErr: false,
		},
		{
			name:    "valid large number",
			input:   "999999999999999999999999999999",
			want:    mustParseBigIntForTest(t, "999999999999999999999999999999"),
			wantErr: false,
		},
		{
			name:    "invalid empty string",
			input:   "",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid non-numeric",
			input:   "abc",
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := BigIntFromString(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want.String(), got.String())
			}
		})
	}
}

func TestBigIntToHex(t *testing.T) {
	tests := []struct {
		name  string
		input *big.Int
		want  string
	}{
		{
			name:  "zero",
			input: big.NewInt(0),
			want:  "0x0",
		},
		{
			name:  "small number",
			input: big.NewInt(255),
			want:  "0xff",
		},
		{
			name:  "large number",
			input: big.NewInt(123456789),
			want:  "0x75bcd15",
		},
		{
			name:  "very large number",
			input: mustParseBigIntForTest(t, "115792089237316195423570985008687907853269984665640564039457584007913129639935"),
			want:  "0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BigIntToHex(tt.input)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestBigIntFromHex(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    *big.Int
		wantErr bool
	}{
		{
			name:    "valid with 0x prefix",
			input:   "0xff",
			want:    big.NewInt(255),
			wantErr: false,
		},
		{
			name:    "valid without 0x prefix",
			input:   "ff",
			want:    big.NewInt(255),
			wantErr: false,
		},
		{
			name:    "valid zero",
			input:   "0x0",
			want:    big.NewInt(0),
			wantErr: false,
		},
		{
			name:    "valid large hex",
			input:   "0x75bcd15",
			want:    big.NewInt(123456789),
			wantErr: false,
		},
		{
			name:    "invalid empty string",
			input:   "",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid non-hex",
			input:   "0xgh",
			want:    nil,
			wantErr: true,
		},
		{
			name:    "invalid with spaces",
			input:   "0x 12",
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := BigIntFromHex(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want.String(), got.String())
			}
		})
	}
}

func mustParseBigIntForTest(t *testing.T, s string) *big.Int {
	val, ok := new(big.Int).SetString(s, 10)
	require.True(t, ok)
	return val
}
