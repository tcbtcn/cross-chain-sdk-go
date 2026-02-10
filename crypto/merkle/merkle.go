package merkle

import (
	"crypto/sha256"
	"sort"
)

type MerkleTree struct {
	leaves [][]byte
	root   []byte
}

func NewMerkleTree(leaves [][]byte) *MerkleTree {
	if len(leaves) == 0 {
		return &MerkleTree{leaves: leaves, root: nil}
	}
	if len(leaves) == 1 {
		return &MerkleTree{leaves: leaves, root: leaves[0]}
	}

	sorted := make([][]byte, len(leaves))
	copy(sorted, leaves)
	sort.Slice(sorted, func(i, j int) bool {
		return string(sorted[i]) < string(sorted[j])
	})

	root := buildTree(sorted)
	return &MerkleTree{leaves: sorted, root: root}
}

func buildTree(leaves [][]byte) []byte {
	if len(leaves) == 1 {
		return leaves[0]
	}

	nextLevel := make([][]byte, 0)
	for i := 0; i < len(leaves); i += 2 {
		if i+1 < len(leaves) {
			left := leaves[i]
			right := leaves[i+1]
			hash := hashPair(left, right)
			nextLevel = append(nextLevel, hash)
		} else {
			nextLevel = append(nextLevel, leaves[i])
		}
	}

	return buildTree(nextLevel)
}

func hashPair(left, right []byte) []byte {
	var combined []byte
	if string(left) < string(right) {
		combined = append(combined, left...)
		combined = append(combined, right...)
	} else {
		combined = append(combined, right...)
		combined = append(combined, left...)
	}
	hash := sha256.Sum256(combined)
	return hash[:]
}

func (mt *MerkleTree) Root() []byte {
	return mt.root
}

func (mt *MerkleTree) GetProof(index int) ([][]byte, error) {
	if index < 0 || index >= len(mt.leaves) {
		return nil, ErrInvalidIndex
	}

	proof := make([][]byte, 0)
	currentLevel := make([][]byte, len(mt.leaves))
	copy(currentLevel, mt.leaves)
	currentIndex := index

	for len(currentLevel) > 1 {
		levelSize := len(currentLevel)
		nextLevel := make([][]byte, 0)
		nextIndex := -1

		for i := 0; i < levelSize; i += 2 {
			if i+1 < levelSize {
				left := currentLevel[i]
				right := currentLevel[i+1]
				hash := hashPair(left, right)
				nextLevel = append(nextLevel, hash)

				if i == currentIndex {
					proof = append(proof, right)
					nextIndex = len(nextLevel) - 1
				} else if i+1 == currentIndex {
					proof = append(proof, left)
					nextIndex = len(nextLevel) - 1
				}
			} else {
				nextLevel = append(nextLevel, currentLevel[i])
				if i == currentIndex {
					nextIndex = len(nextLevel) - 1
				}
			}
		}

		currentLevel = nextLevel
		if nextIndex >= 0 {
			currentIndex = nextIndex
		}
	}

	return proof, nil
}

var ErrInvalidIndex = &MerkleError{Message: "invalid index"}

type MerkleError struct {
	Message string
}

func (e *MerkleError) Error() string {
	return e.Message
}

func VerifyProof(leaf []byte, proof [][]byte, root []byte) bool {
	if len(proof) == 0 {
		return false
	}

	current := leaf

	for _, proofNode := range proof {
		current = hashPair(current, proofNode)
	}

	return string(current) == string(root)
}
