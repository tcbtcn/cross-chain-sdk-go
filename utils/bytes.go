package utils

import (
	"encoding/hex"
	"fmt"
)

func BufferFromHex(hexStr string, bytesSize int) ([]byte, error) {
	if !IsHexBytes(hexStr) {
		return nil, fmt.Errorf("invalid hex string: %s", hexStr)
	}

	hexBytes := hexStr
	if len(hexBytes) >= 2 && hexBytes[0:2] == "0x" {
		hexBytes = hexBytes[2:]
	}

	data, err := hex.DecodeString(hexBytes)
	if err != nil {
		return nil, err
	}

	if bytesSize > 0 {
		if len(data) > bytesSize {
			return nil, fmt.Errorf("buffer size exceeds %d bytes", bytesSize)
		}
		padded := make([]byte, bytesSize)
		copy(padded[bytesSize-len(data):], data)
		return padded, nil
	}

	return data, nil
}

func BufferToHex(buf []byte) string {
	return "0x" + hex.EncodeToString(buf)
}

func IsHexBytes(s string) bool {
	if len(s) < 2 {
		return false
	}
	if s[0:2] != "0x" {
		return false
	}
	hexPart := s[2:]
	if len(hexPart)%2 != 0 {
		return false
	}
	_, err := hex.DecodeString(hexPart)
	return err == nil
}

func GetBytesCount(hexStr string) (int64, error) {
	if !IsHexBytes(hexStr) {
		return 0, fmt.Errorf("invalid hex string")
	}
	hexPart := hexStr[2:]
	return int64(len(hexPart) / 2), nil
}
