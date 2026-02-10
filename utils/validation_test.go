package utils

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateUInteger(t *testing.T) {
	tests := []struct {
		name    string
		val     *big.Int
		max     *big.Int
		wantErr bool
	}{
		{
			name:    "valid positive number",
			val:     big.NewInt(100),
			max:     big.NewInt(1000),
			wantErr: false,
		},
		{
			name:    "valid zero",
			val:     big.NewInt(0),
			max:     big.NewInt(1000),
			wantErr: false,
		},
		{
			name:    "valid at max",
			val:     big.NewInt(1000),
			max:     big.NewInt(1000),
			wantErr: false,
		},
		{
			name:    "invalid negative",
			val:     big.NewInt(-1),
			max:     big.NewInt(1000),
			wantErr: true,
		},
		{
			name:    "invalid exceeds max",
			val:     big.NewInt(1001),
			max:     big.NewInt(1000),
			wantErr: true,
		},
		{
			name:    "valid with nil max",
			val:     big.NewInt(1000000),
			max:     nil,
			wantErr: false,
		},
		{
			name:    "invalid negative with nil max",
			val:     big.NewInt(-1),
			max:     nil,
			wantErr: true,
		},
		{
			name:    "valid very large number",
			val:     mustParseBigIntForValidationTest(t, "999999999999999999999999999999"),
			max:     mustParseBigIntForValidationTest(t, "1000000000000000000000000000000"),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateUInteger(tt.val, tt.max)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestValidateHexBytes(t *testing.T) {
	tests := []struct {
		name         string
		hexStr       string
		expectedSize int64
		wantErr      bool
	}{
		{
			name:         "valid 32 bytes",
			hexStr:       repeatHex("00", 32),
			expectedSize: 32,
			wantErr:      false,
		},
		{
			name:         "valid 16 bytes",
			hexStr:       repeatHex("00", 16),
			expectedSize: 16,
			wantErr:      false,
		},
		{
			name:         "valid 1 byte",
			hexStr:       "0x00",
			expectedSize: 1,
			wantErr:      false,
		},
		{
			name:         "invalid wrong size",
			hexStr:       repeatHex("00", 16),
			expectedSize: 32,
			wantErr:      true,
		},
		{
			name:         "invalid not hex",
			hexStr:       "0xgh",
			expectedSize: 1,
			wantErr:      true,
		},
		{
			name:         "invalid without 0x",
			hexStr:       "00",
			expectedSize: 1,
			wantErr:      true,
		},
		{
			name:         "invalid empty",
			hexStr:       "",
			expectedSize: 0,
			wantErr:      true,
		},
		{
			name:         "invalid odd length",
			hexStr:       "0x123",
			expectedSize: 1,
			wantErr:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateHexBytes(tt.hexStr, tt.expectedSize)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func mustParseBigIntForValidationTest(t *testing.T, s string) *big.Int {
	val, ok := new(big.Int).SetString(s, 10)
	if !ok {
		t.Fatalf("failed to parse big.Int: %s", s)
	}
	return val
}
