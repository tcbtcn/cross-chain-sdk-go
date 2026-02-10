package hashlock

import (
	"math/big"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func repeatHex(s string, n int) string {
	return "0x" + strings.Repeat(s, n)
}

var (
	sampleSecret1 = repeatHex("00", 31) + "01"
	sampleSecret2 = repeatHex("00", 31) + "02"
	sampleSecret3 = repeatHex("00", 31) + "03"
	sampleSecret4 = repeatHex("00", 31) + "04"
)

func TestHashSecret(t *testing.T) {
	secret := sampleSecret1
	hash, err := HashSecret(secret)
	require.NoError(t, err)
	assert.True(t, len(hash) > 2)
	assert.Equal(t, "0x", hash[:2])

	// Hash should be deterministic
	hash2, err := HashSecret(secret)
	require.NoError(t, err)
	assert.Equal(t, hash, hash2)
}

func TestHashSecret_Invalid(t *testing.T) {
	tests := []struct {
		name    string
		secret  string
		wantErr bool
	}{
		{
			name:    "invalid length",
			secret:  "0x123",
			wantErr: true,
		},
		{
			name:    "not hex",
			secret:  "nothex",
			wantErr: true,
		},
		{
			name:    "empty",
			secret:  "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := HashSecret(tt.secret)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Empty(t, hash)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestForSingleFill(t *testing.T) {
	secret := sampleSecret1
	hashLock, err := ForSingleFill(secret)
	require.NoError(t, err)
	assert.NotNil(t, hashLock)

	// GetPartsCount is only meaningful for multiple fills
	// For single fill, it returns whatever is in bits 240-247 of the hash
	partsCount, err := hashLock.GetPartsCount()
	require.NoError(t, err)
	assert.NotNil(t, partsCount)
}

func TestForMultipleFills(t *testing.T) {
	secrets := []string{
		sampleSecret1,
		sampleSecret2,
		sampleSecret3,
		sampleSecret4,
	}

	leaves, err := GetMerkleLeaves(secrets)
	require.NoError(t, err)
	assert.Equal(t, len(secrets), len(leaves))

	hashLock, err := ForMultipleFills(leaves)
	require.NoError(t, err)
	assert.NotNil(t, hashLock)

	partsCount, err := hashLock.GetPartsCount()
	require.NoError(t, err)
	assert.Equal(t, big.NewInt(int64(len(secrets)-1)).String(), partsCount.String())
}

func TestForMultipleFills_Invalid(t *testing.T) {
	tests := []struct {
		name    string
		leaves  []string
		wantErr bool
	}{
		{
			name:    "too few leaves",
			leaves:  []string{repeatHex("00", 32)},
			wantErr: true,
		},
		{
			name:    "exactly 2 leaves",
			leaves:  []string{repeatHex("00", 32), repeatHex("00", 32)},
			wantErr: true,
		},
		{
			name:    "empty leaves",
			leaves:  []string{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hashLock, err := ForMultipleFills(tt.leaves)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, hashLock)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetMerkleLeaves(t *testing.T) {
	secrets := []string{
		sampleSecret1,
		sampleSecret2,
		sampleSecret3,
	}

	leaves, err := GetMerkleLeaves(secrets)
	require.NoError(t, err)
	assert.Equal(t, len(secrets), len(leaves))

	for _, leaf := range leaves {
		assert.True(t, len(leaf) > 2, "leaf should have at least 2 characters")
		if len(leaf) >= 2 {
			assert.Equal(t, "0x", leaf[:2])
		}
	}
}

func TestHashLock_GetPartsCount(t *testing.T) {
	// Single fill - GetPartsCount is undefined for single fill
	// (returns whatever is in bits 240-247 of the hash)
	secret := sampleSecret1
	singleHashLock, err := ForSingleFill(secret)
	require.NoError(t, err)

	partsCount, err := singleHashLock.GetPartsCount()
	require.NoError(t, err)
	assert.NotNil(t, partsCount)

	// Multiple fills
	secrets := []string{sampleSecret1, sampleSecret2, sampleSecret3, sampleSecret4}
	leaves, err := GetMerkleLeaves(secrets)
	require.NoError(t, err)
	multipleHashLock, err := ForMultipleFills(leaves)
	require.NoError(t, err)

	partsCount, err = multipleHashLock.GetPartsCount()
	require.NoError(t, err)
	assert.Equal(t, big.NewInt(3).String(), partsCount.String())
}

func TestHashLock_ToString(t *testing.T) {
	secret := sampleSecret1
	hashLock, err := ForSingleFill(secret)
	require.NoError(t, err)

	str := hashLock.ToString()
	assert.True(t, len(str) > 2)
	assert.Equal(t, "0x", str[:2])
}

func TestHashLock_ToBuffer(t *testing.T) {
	secret := sampleSecret1
	hashLock, err := ForSingleFill(secret)
	require.NoError(t, err)

	buf, err := hashLock.ToBuffer()
	require.NoError(t, err)
	assert.Equal(t, 32, len(buf))
}

func TestHashLock_Equal(t *testing.T) {
	secret := sampleSecret1
	hashLock1, err := ForSingleFill(secret)
	require.NoError(t, err)

	hashLock2, err := ForSingleFill(secret)
	require.NoError(t, err)

	assert.True(t, hashLock1.Equal(hashLock2))

	secret2 := sampleSecret2
	hashLock3, err := ForSingleFill(secret2)
	require.NoError(t, err)

	assert.False(t, hashLock1.Equal(hashLock3))
}

func TestNewHashLock(t *testing.T) {
	validHash := repeatHex("00", 32)
	hashLock, err := NewHashLock(validHash)
	require.NoError(t, err)
	assert.NotNil(t, hashLock)

	invalidHash := "0x123"
	hashLock2, err := NewHashLock(invalidHash)
	assert.Error(t, err)
	assert.Nil(t, hashLock2)
}
