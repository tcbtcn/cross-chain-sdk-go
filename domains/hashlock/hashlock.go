package hashlock

import (
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/dawitel/cross-chain-sdk-go/crypto"
	"github.com/dawitel/cross-chain-sdk-go/crypto/merkle"
	"github.com/dawitel/cross-chain-sdk-go/utils"
)

const Web3Type = "bytes32"

type HashLock struct {
	value string
}

func NewHashLock(value string) (*HashLock, error) {
	if err := utils.ValidateHexBytes(value, 32); err != nil {
		return nil, fmt.Errorf("HashLock value must be bytes32 hex encoded: %w", err)
	}
	return &HashLock{value: value}, nil
}

func HashLockFromString(value string) (*HashLock, error) {
	return NewHashLock(value)
}

func HashLockFromBuffer(buf []byte) (*HashLock, error) {
	hexStr := "0x" + hex.EncodeToString(buf)
	return HashLockFromString(hexStr)
}

func HashSecret(secret string) (string, error) {
	if err := utils.ValidateHexBytes(secret, 32); err != nil {
		return "", fmt.Errorf("secret length must be 32 bytes hex encoded: %w", err)
	}
	secretBytes, _ := utils.BufferFromHex(secret, -1)
	hash := crypto.Keccak256(secretBytes)
	return "0x" + hex.EncodeToString(hash), nil
}

func GetMerkleLeaves(secrets []string) ([]string, error) {
	secretHashes := make([]string, len(secrets))
	for i, secret := range secrets {
		hash, err := HashSecret(secret)
		if err != nil {
			return nil, err
		}
		secretHashes[i] = hash
	}
	return GetMerkleLeavesFromSecretHashes(secretHashes)
}

func GetMerkleLeavesFromSecretHashes(secretHashes []string) ([]string, error) {
	leaves := make([]string, len(secretHashes))
	for i, hash := range secretHashes {
		hashBytes, err := utils.BufferFromHex(hash, 32)
		if err != nil {
			return nil, err
		}
		packedHash, err := crypto.SolidityPackedKeccak256(
			[]string{"uint64", "bytes32"},
			[]interface{}{uint64(i), hashBytes},
		)
		if err != nil {
			return nil, err
		}
		leaves[i] = "0x" + hex.EncodeToString(packedHash)
	}
	return leaves, nil
}

func ForSingleFill(secret string) (*HashLock, error) {
	hash, err := HashSecret(secret)
	if err != nil {
		return nil, err
	}
	return NewHashLock(hash)
}

func ForMultipleFills(leaves []string) (*HashLock, error) {
	if len(leaves) <= 2 {
		return nil, fmt.Errorf("leaves array must be greater than 2. Or use ForSingleFill")
	}

	leafBytes := make([][]byte, len(leaves))
	for i, leaf := range leaves {
		bytes, err := utils.BufferFromHex(leaf, 32)
		if err != nil {
			return nil, err
		}
		leafBytes[i] = bytes
	}

	tree := merkle.NewMerkleTree(leafBytes)
	root := tree.Root()

	rootBig := new(big.Int).SetBytes(root)
	count := big.NewInt(int64(len(leaves) - 1))

	// Clear bits 240-255 (16 bits) first
	mask := new(big.Int)
	for i := 0; i < 16; i++ {
		mask.SetBit(mask, 240+i, 1)
	}
	rootBig.AndNot(rootBig, mask)

	for i := 0; i < 16 && i < count.BitLen(); i++ {
		if count.Bit(i) == 1 {
			rootBig.SetBit(rootBig, 240+i, 1)
		}
	}

	rootHex := "0x" + fmt.Sprintf("%064x", rootBig)
	return NewHashLock(rootHex)
}

func (hl *HashLock) GetPartsCount() (*big.Int, error) {
	val, err := utils.BigIntFromHex(hl.value)
	if err != nil {
		return nil, err
	}
	// Extract bits 240-255 (16 bits) as per TypeScript SDK BitMask(240n, 256n)
	mask := new(big.Int)
	for i := 0; i < 16; i++ {
		mask.SetBit(mask, 240+i, 1)
	}
	result := new(big.Int).And(val, mask)
	result.Rsh(result, 240)
	return result, nil
}

func (hl *HashLock) ToString() string {
	return hl.value
}

func (hl *HashLock) ToJSON() string {
	return hl.value
}

func (hl *HashLock) ToBuffer() ([]byte, error) {
	return utils.BufferFromHex(hl.value, 32)
}

func (hl *HashLock) Equal(other *HashLock) bool {
	return hl.value == other.value
}
