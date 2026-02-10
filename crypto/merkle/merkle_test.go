package merkle

import (
	"crypto/sha256"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewMerkleTree(t *testing.T) {
	tests := []struct {
		name   string
		leaves [][]byte
	}{
		{
			name:   "empty leaves",
			leaves: [][]byte{},
		},
		{
			name:   "single leaf",
			leaves: [][]byte{[]byte("leaf1")},
		},
		{
			name:   "two leaves",
			leaves: [][]byte{[]byte("leaf1"), []byte("leaf2")},
		},
		{
			name:   "four leaves",
			leaves: [][]byte{[]byte("leaf1"), []byte("leaf2"), []byte("leaf3"), []byte("leaf4")},
		},
		{
			name: "eight leaves",
			leaves: [][]byte{
				[]byte("leaf1"), []byte("leaf2"), []byte("leaf3"), []byte("leaf4"),
				[]byte("leaf5"), []byte("leaf6"), []byte("leaf7"), []byte("leaf8"),
			},
		},
		{
			name:   "three leaves (odd number)",
			leaves: [][]byte{[]byte("leaf1"), []byte("leaf2"), []byte("leaf3")},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tree := NewMerkleTree(tt.leaves)
			require.NotNil(t, tree)

			if len(tt.leaves) == 0 {
				assert.Nil(t, tree.Root())
			} else if len(tt.leaves) == 1 {
				assert.Equal(t, tt.leaves[0], tree.Root())
			} else {
				assert.NotNil(t, tree.Root())
				assert.Equal(t, 32, len(tree.Root()))
			}
		})
	}
}

func TestMerkleTree_Root(t *testing.T) {
	leaves := [][]byte{
		[]byte("leaf1"),
		[]byte("leaf2"),
		[]byte("leaf3"),
		[]byte("leaf4"),
	}

	tree := NewMerkleTree(leaves)
	root := tree.Root()

	assert.NotNil(t, root)
	assert.Equal(t, 32, len(root))
}

func TestMerkleTree_GetProof(t *testing.T) {
	leaves := [][]byte{
		[]byte("leaf1"),
		[]byte("leaf2"),
		[]byte("leaf3"),
		[]byte("leaf4"),
	}

	tree := NewMerkleTree(leaves)

	for i := range leaves {
		proof, err := tree.GetProof(i)
		require.NoError(t, err)
		assert.NotNil(t, proof)
	}
}

func TestMerkleTree_VerifyProof(t *testing.T) {
	leaves := [][]byte{
		[]byte("leaf1"),
		[]byte("leaf2"),
		[]byte("leaf3"),
		[]byte("leaf4"),
	}

	tree := NewMerkleTree(leaves)
	root := tree.Root()

	for i, leaf := range leaves {
		proof, err := tree.GetProof(i)
		require.NoError(t, err)
		verified := VerifyProof(leaf, proof, root)
		assert.True(t, verified, "proof should verify for leaf %d", i)
	}
}

func TestMerkleTree_VerifyProof_Invalid(t *testing.T) {
	leaves := [][]byte{
		[]byte("leaf1"),
		[]byte("leaf2"),
		[]byte("leaf3"),
		[]byte("leaf4"),
	}

	tree := NewMerkleTree(leaves)
	root := tree.Root()

	// Try to verify wrong leaf
	wrongLeaf := []byte("wrong")
	proof, err := tree.GetProof(0)
	require.NoError(t, err)
	verified := VerifyProof(wrongLeaf, proof, root)
	assert.False(t, verified, "proof should not verify for wrong leaf")
}

func TestMerkleTree_Deterministic(t *testing.T) {
	leaves := [][]byte{
		[]byte("leaf1"),
		[]byte("leaf2"),
		[]byte("leaf3"),
		[]byte("leaf4"),
	}

	tree1 := NewMerkleTree(leaves)
	tree2 := NewMerkleTree(leaves)

	assert.Equal(t, tree1.Root(), tree2.Root(), "trees with same leaves should have same root")
}

func TestMerkleTree_Sorted(t *testing.T) {
	leaves1 := [][]byte{
		[]byte("leaf1"),
		[]byte("leaf2"),
		[]byte("leaf3"),
		[]byte("leaf4"),
	}

	leaves2 := [][]byte{
		[]byte("leaf4"),
		[]byte("leaf3"),
		[]byte("leaf2"),
		[]byte("leaf1"),
	}

	tree1 := NewMerkleTree(leaves1)
	tree2 := NewMerkleTree(leaves2)

	assert.Equal(t, tree1.Root(), tree2.Root(), "trees should sort leaves before building")
}

func TestHashPair(t *testing.T) {
	left := []byte("left")
	right := []byte("right")

	hash := hashPair(left, right)

	assert.Equal(t, 32, len(hash))

	// Verify it's SHA256 of sorted concatenated raw bytes
	var combined []byte
	if string(left) < string(right) {
		combined = append(combined, left...)
		combined = append(combined, right...)
	} else {
		combined = append(combined, right...)
		combined = append(combined, left...)
	}
	expected := sha256.Sum256(combined)

	assert.Equal(t, expected[:], hash)
}
