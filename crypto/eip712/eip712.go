package eip712

import (
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
)

type TypedData struct {
	Types       map[string]interface{} `json:"types"`
	PrimaryType string                 `json:"primaryType"`
	Domain      map[string]interface{} `json:"domain"`
	Message     map[string]interface{} `json:"message"`
}

func HashTypedData(typedData TypedData) ([]byte, error) {
	domainSeparator, err := typedData.HashStruct("EIP712Domain", typedData.Domain)
	if err != nil {
		return nil, err
	}

	typedDataHash, err := typedData.HashStruct(typedData.PrimaryType, typedData.Message)
	if err != nil {
		return nil, err
	}

	rawData := []byte(fmt.Sprintf("\x19\x01%s%s", string(domainSeparator), string(typedDataHash)))
	return crypto.Keccak256(rawData), nil
}

func (td *TypedData) HashStruct(primaryType string, data map[string]interface{}) ([]byte, error) {
	encodedData, err := td.EncodeData(primaryType, data)
	if err != nil {
		return nil, err
	}
	return crypto.Keccak256(encodedData), nil
}

func (td *TypedData) EncodeData(primaryType string, data map[string]interface{}) ([]byte, error) {
	typesList, ok := td.Types[primaryType].([]interface{})
	if !ok {
		return nil, fmt.Errorf("types not found for %s", primaryType)
	}

	var encoded []byte

	for _, fieldInterface := range typesList {
		field, ok := fieldInterface.(map[string]interface{})
		if !ok {
			continue
		}
		fieldName, _ := field["name"].(string)
		fieldType, _ := field["type"].(string)
		value := data[fieldName]
		encodedValue, err := td.EncodeField(fieldType, value)
		if err != nil {
			return nil, err
		}
		encoded = append(encoded, encodedValue...)
	}

	return encoded, nil
}

func (td *TypedData) EncodeField(fieldType string, value interface{}) ([]byte, error) {
	if fieldType == "string" {
		str, ok := value.(string)
		if !ok {
			return nil, fmt.Errorf("expected string, got %T", value)
		}
		return crypto.Keccak256([]byte(str)), nil
	}

	if fieldType == "bytes" {
		bytes, ok := value.([]byte)
		if !ok {
			return nil, fmt.Errorf("expected []byte, got %T", value)
		}
		return crypto.Keccak256(bytes), nil
	}

	if fieldType == "bytes32" {
		var bytes []byte
		switch v := value.(type) {
		case string:
			if len(v) >= 2 && v[0:2] == "0x" {
				bytes, _ = hex.DecodeString(v[2:])
			} else {
				bytes, _ = hex.DecodeString(v)
			}
		case []byte:
			bytes = v
		default:
			return nil, fmt.Errorf("expected bytes32, got %T", value)
		}
		if len(bytes) < 32 {
			padded := make([]byte, 32)
			copy(padded[32-len(bytes):], bytes)
			bytes = padded
		}
		return bytes[:32], nil
	}

	if fieldType == "address" {
		var addr common.Address
		switch v := value.(type) {
		case string:
			addr = common.HexToAddress(v)
		case common.Address:
			addr = v
		default:
			return nil, fmt.Errorf("expected address, got %T", value)
		}
		return addr.Bytes(), nil
	}

	if fieldType == "uint256" || fieldType == "uint128" || fieldType == "uint64" || fieldType == "uint32" {
		var num *big.Int
		switch v := value.(type) {
		case string:
			var ok bool
			num, ok = new(big.Int).SetString(v, 10)
			if !ok {
				return nil, fmt.Errorf("invalid number string: %s", v)
			}
		case *big.Int:
			num = v
		case int64:
			num = big.NewInt(v)
		case uint64:
			num = new(big.Int).SetUint64(v)
		default:
			return nil, fmt.Errorf("expected number, got %T", value)
		}
		return math.U256Bytes(num), nil
	}

	if len(fieldType) >= 2 && fieldType[len(fieldType)-2:] == "[]" {
		baseType := fieldType[:len(fieldType)-2]
		slice, ok := value.([]interface{})
		if !ok {
			return nil, fmt.Errorf("expected array, got %T", value)
		}
		encoded := crypto.Keccak256([]byte(fmt.Sprintf("%d", len(slice))))
		for _, item := range slice {
			itemEncoded, err := td.EncodeField(baseType, item)
			if err != nil {
				return nil, err
			}
			encoded = append(encoded, itemEncoded...)
		}
		return crypto.Keccak256(encoded), nil
	}

	if _, exists := td.Types[fieldType]; exists {
		nestedData, ok := value.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("expected struct, got %T", value)
		}
		return td.HashStruct(fieldType, nestedData)
	}

	return nil, fmt.Errorf("unsupported type: %s", fieldType)
}

func SignTypedData(typedData TypedData, privateKey []byte) ([]byte, error) {
	hash, err := HashTypedData(typedData)
	if err != nil {
		return nil, err
	}

	privateKeyECDSA, err := crypto.ToECDSA(privateKey)
	if err != nil {
		return nil, err
	}

	return crypto.Sign(hash, privateKeyECDSA)
}

func RecoverTypedData(typedData TypedData, signature []byte) (common.Address, error) {
	hash, err := HashTypedData(typedData)
	if err != nil {
		return common.Address{}, err
	}

	if len(signature) != 65 {
		return common.Address{}, fmt.Errorf("invalid signature length: %d", len(signature))
	}

	if signature[64] >= 27 {
		signature[64] -= 27
	}

	pubkey, err := crypto.SigToPub(hash, signature)
	if err != nil {
		return common.Address{}, err
	}

	return crypto.PubkeyToAddress(*pubkey), nil
}

func FormatSignature(sig []byte) string {
	if len(sig) != 65 {
		return ""
	}
	if sig[64] >= 27 {
		sig[64] -= 27
	}
	return "0x" + hex.EncodeToString(sig)
}
