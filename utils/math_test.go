package utils

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMulDiv(t *testing.T) {
	tests := []struct {
		name        string
		a           *big.Int
		b           *big.Int
		denominator *big.Int
		want        *big.Int
	}{
		{
			name:        "simple multiplication and division",
			a:           big.NewInt(10),
			b:           big.NewInt(20),
			denominator: big.NewInt(2),
			want:        big.NewInt(100),
		},
		{
			name:        "zero numerator",
			a:           big.NewInt(0),
			b:           big.NewInt(20),
			denominator: big.NewInt(2),
			want:        big.NewInt(0),
		},
		{
			name:        "zero denominator result",
			a:           big.NewInt(10),
			b:           big.NewInt(0),
			denominator: big.NewInt(2),
			want:        big.NewInt(0),
		},
		{
			name:        "large numbers",
			a:           big.NewInt(1000000),
			b:           big.NewInt(2000000),
			denominator: big.NewInt(1000),
			want:        big.NewInt(2000000000),
		},
		{
			name:        "division with remainder",
			a:           big.NewInt(7),
			b:           big.NewInt(3),
			denominator: big.NewInt(2),
			want:        big.NewInt(10), // (7 * 3) / 2 = 21 / 2 = 10 (truncated)
		},
		{
			name:        "very large numbers",
			a:           mustParseBigIntForMathTest(t, "1000000000000000000000000"),
			b:           mustParseBigIntForMathTest(t, "2000000000000000000000000"),
			denominator: big.NewInt(1000000),
			want:        mustParseBigIntForMathTest(t, "2000000000000000000000000000000000000000000"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := MulDiv(tt.a, tt.b, tt.denominator)
			assert.Equal(t, tt.want.String(), got.String())
		})
	}
}

func TestMulDiv_DoesNotModifyInputs(t *testing.T) {
	a := big.NewInt(10)
	b := big.NewInt(20)
	denominator := big.NewInt(2)

	aOriginal := new(big.Int).Set(a)
	bOriginal := new(big.Int).Set(b)
	denominatorOriginal := new(big.Int).Set(denominator)

	_ = MulDiv(a, b, denominator)

	assert.Equal(t, aOriginal.String(), a.String(), "a should not be modified")
	assert.Equal(t, bOriginal.String(), b.String(), "b should not be modified")
	assert.Equal(t, denominatorOriginal.String(), denominator.String(), "denominator should not be modified")
}

func mustParseBigIntForMathTest(t *testing.T, s string) *big.Int {
	val, ok := new(big.Int).SetString(s, 10)
	if !ok {
		t.Fatalf("failed to parse big.Int: %s", s)
	}
	return val
}
