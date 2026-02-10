package utils

import (
	"fmt"
	"math/big"
	"strings"
)

func IsBigIntString(s string) bool {
	_, ok := new(big.Int).SetString(s, 10)
	return ok
}

func BigIntFromString(s string) (*big.Int, error) {
	val, ok := new(big.Int).SetString(s, 10)
	if !ok {
		return nil, fmt.Errorf("invalid bigint string: %s", s)
	}
	return val, nil
}

func BigIntToHex(val *big.Int) string {
	hex := val.Text(16)
	if hex == "0" {
		return "0x0"
	}
	return "0x" + hex
}

func BigIntFromHex(hexStr string) (*big.Int, error) {
	hexStr = strings.TrimPrefix(hexStr, "0x")
	val, ok := new(big.Int).SetString(hexStr, 16)
	if !ok {
		return nil, fmt.Errorf("invalid hex string: %s", hexStr)
	}
	return val, nil
}
