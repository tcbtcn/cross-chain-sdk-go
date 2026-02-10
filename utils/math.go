package utils

import "math/big"

func MulDiv(a, b, denominator *big.Int) *big.Int {
	result := new(big.Int).Mul(a, b)
	result.Div(result, denominator)
	return result
}
