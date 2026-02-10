package utils

import (
	"fmt"
	"math/big"
)

func ValidateUInteger(val *big.Int, max *big.Int) error {
	if val.Sign() < 0 {
		return fmt.Errorf("value must be non-negative")
	}
	if max != nil && val.Cmp(max) > 0 {
		return fmt.Errorf("value exceeds maximum: %s", max.String())
	}
	return nil
}

func ValidateHexBytes(hexStr string, expectedSize int64) error {
	if !IsHexBytes(hexStr) {
		return fmt.Errorf("invalid hex string")
	}
	count, err := GetBytesCount(hexStr)
	if err != nil {
		return err
	}
	if count != expectedSize {
		return fmt.Errorf("expected %d bytes, got %d", expectedSize, count)
	}
	return nil
}
