package crypto

import (
	"encoding/hex"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func EncodeCancelOrder(orderHash string, makerTraits *big.Int) (string, error) {
	orderHashBytes := common.HexToHash(orderHash).Bytes()

	methodID := crypto.Keccak256([]byte("cancelOrder(uint256,bytes32)"))[:4]

	makerTraitsBytes := make([]byte, 32)
	makerTraits.FillBytes(makerTraitsBytes)

	callData := append(methodID, makerTraitsBytes...)
	callData = append(callData, orderHashBytes...)

	return "0x" + hex.EncodeToString(callData), nil
}
