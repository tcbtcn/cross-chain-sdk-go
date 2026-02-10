package crypto

import (
	"encoding/hex"
	"fmt"
	"math/big"

	"golang.org/x/crypto/sha3"
)

func Keccak256(data []byte) []byte {
	hash := sha3.NewLegacyKeccak256()
	hash.Write(data)
	return hash.Sum(nil)
}

func Keccak256Hex(data []byte) string {
	return "0x" + hex.EncodeToString(Keccak256(data))
}

func Keccak256FromHex(hexStr string) ([]byte, error) {
	if len(hexStr) < 2 || hexStr[0:2] != "0x" {
		return nil, fmt.Errorf("hex string must start with 0x")
	}
	hexStr = hexStr[2:]
	data, err := hex.DecodeString(hexStr)
	if err != nil {
		return nil, err
	}
	return Keccak256(data), nil
}

func SolidityPackedKeccak256(types []string, values []interface{}) ([]byte, error) {
	data, err := packSolidityData(types, values)
	if err != nil {
		return nil, err
	}
	return Keccak256(data), nil
}

func SolidityPackedKeccak256Hex(types []string, values []interface{}) (string, error) {
	hash, err := SolidityPackedKeccak256(types, values)
	if err != nil {
		return "", err
	}
	return "0x" + hex.EncodeToString(hash), nil
}

func packSolidityData(types []string, values []interface{}) ([]byte, error) {
	if len(types) != len(values) {
		return nil, &PackingError{Message: "types and values length mismatch"}
	}
	var result []byte
	for i, typ := range types {
		val := values[i]
		switch typ {
		case "uint64":
			var num uint64
			switch v := val.(type) {
			case uint64:
				num = v
			case int:
				num = uint64(v)
			case int64:
				num = uint64(v)
			default:
				return nil, &PackingError{Message: "invalid uint64 value"}
			}
			buf := make([]byte, 8)
			big.NewInt(int64(num)).FillBytes(buf)
			result = append(result, buf...)
		case "bytes32":
			var data []byte
			switch v := val.(type) {
			case string:
				if len(v) >= 2 && v[0:2] == "0x" {
					data, _ = hex.DecodeString(v[2:])
				} else {
					data, _ = hex.DecodeString(v)
				}
				if len(data) < 32 {
					padded := make([]byte, 32)
					copy(padded[32-len(data):], data)
					data = padded
				} else if len(data) > 32 {
					data = data[:32]
				}
			case []byte:
				if len(v) < 32 {
					padded := make([]byte, 32)
					copy(padded[32-len(v):], v)
					data = padded
				} else if len(v) > 32 {
					data = v[:32]
				} else {
					data = v
				}
			default:
				return nil, &PackingError{Message: "invalid bytes32 value"}
			}
			result = append(result, data...)
		default:
			return nil, &PackingError{Message: "unsupported type: " + typ}
		}
	}
	return result, nil
}

type PackingError struct {
	Message string
}

func (e *PackingError) Error() string {
	return e.Message
}
