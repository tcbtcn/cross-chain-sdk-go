package timelocks

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTimeLocks(t *testing.T) {
	params := TimeLocksParams{
		SrcWithdrawal:         big.NewInt(3600),
		SrcPublicWithdrawal:   big.NewInt(7200),
		SrcCancellation:       big.NewInt(10800),
		SrcPublicCancellation: big.NewInt(14400),
		DstWithdrawal:         big.NewInt(3600),
		DstPublicWithdrawal:   big.NewInt(7200),
		DstCancellation:       big.NewInt(10800),
		DeployedAt:            big.NewInt(0),
	}

	tl, err := NewTimeLocks(params)
	require.NoError(t, err)
	assert.NotNil(t, tl)
}

func TestNewTimeLocks_Validation(t *testing.T) {
	tests := []struct {
		name    string
		params  TimeLocksParams
		wantErr bool
	}{
		{
			name: "valid params",
			params: TimeLocksParams{
				SrcWithdrawal:         big.NewInt(3600),
				SrcPublicWithdrawal:   big.NewInt(7200),
				SrcCancellation:       big.NewInt(10800),
				SrcPublicCancellation: big.NewInt(14400),
				DstWithdrawal:         big.NewInt(3600),
				DstPublicWithdrawal:   big.NewInt(7200),
				DstCancellation:       big.NewInt(10800),
			},
			wantErr: false,
		},
		{
			name: "invalid srcWithdrawal >= srcPublicWithdrawal",
			params: TimeLocksParams{
				SrcWithdrawal:       big.NewInt(7200),
				SrcPublicWithdrawal: big.NewInt(3600),
			},
			wantErr: true,
		},
		{
			name: "invalid exceeds uint32 max",
			params: TimeLocksParams{
				SrcWithdrawal:         mustParseBigInt(t, "4294967296"), // > uint32 max
				SrcPublicWithdrawal:   big.NewInt(7200),
				SrcCancellation:       big.NewInt(1800),
				SrcPublicCancellation: big.NewInt(3600),
				DstWithdrawal:         big.NewInt(3600),
				DstPublicWithdrawal:   big.NewInt(7200),
				DstCancellation:       big.NewInt(1800),
			},
			wantErr: true,
		},
		{
			name: "negative value",
			params: TimeLocksParams{
				SrcWithdrawal: big.NewInt(-1),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tl, err := NewTimeLocks(tt.params)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, tl)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, tl)
			}
		})
	}
}

func TestTimeLocks_Build(t *testing.T) {
	params := TimeLocksParams{
		SrcWithdrawal:         big.NewInt(3600),
		SrcPublicWithdrawal:   big.NewInt(7200),
		SrcCancellation:       big.NewInt(10800),
		SrcPublicCancellation: big.NewInt(14400),
		DstWithdrawal:         big.NewInt(3600),
		DstPublicWithdrawal:   big.NewInt(7200),
		DstCancellation:       big.NewInt(10800),
		DeployedAt:            big.NewInt(0),
	}

	tl, err := NewTimeLocks(params)
	require.NoError(t, err)

	built := tl.Build()
	assert.NotNil(t, built)
}

func TestTimeLocks_DeployedAt(t *testing.T) {
	params := TimeLocksParams{
		SrcWithdrawal:         big.NewInt(3600),
		SrcPublicWithdrawal:   big.NewInt(7200),
		SrcCancellation:       big.NewInt(10800),
		SrcPublicCancellation: big.NewInt(14400),
		DstWithdrawal:         big.NewInt(3600),
		DstPublicWithdrawal:   big.NewInt(7200),
		DstCancellation:       big.NewInt(10800),
		DeployedAt:            big.NewInt(1000),
	}

	tl, err := NewTimeLocks(params)
	require.NoError(t, err)

	deployedAt := tl.DeployedAt()
	assert.Equal(t, big.NewInt(1000).String(), deployedAt.String())
}

func TestTimeLocks_SetDeployedAt(t *testing.T) {
	params := TimeLocksParams{
		SrcWithdrawal:         big.NewInt(3600),
		SrcPublicWithdrawal:   big.NewInt(7200),
		SrcCancellation:       big.NewInt(10800),
		SrcPublicCancellation: big.NewInt(14400),
		DstWithdrawal:         big.NewInt(3600),
		DstPublicWithdrawal:   big.NewInt(7200),
		DstCancellation:       big.NewInt(10800),
	}

	tl, err := NewTimeLocks(params)
	require.NoError(t, err)

	newTime := big.NewInt(2000)
	result := tl.SetDeployedAt(newTime)
	assert.Equal(t, tl, result)
	assert.Equal(t, newTime.String(), tl.DeployedAt().String())
}

func TestTimeLocksFromBigInt(t *testing.T) {
	// Create a valid bigint from time locks
	// Note: The order matters - when decoding from bigint, params are in reverse order
	params := TimeLocksParams{
		SrcWithdrawal:         big.NewInt(3600),
		SrcPublicWithdrawal:   big.NewInt(7200),
		SrcCancellation:       big.NewInt(10800),
		SrcPublicCancellation: big.NewInt(14400),
		DstWithdrawal:         big.NewInt(3600),
		DstPublicWithdrawal:   big.NewInt(7200),
		DstCancellation:       big.NewInt(10800),
		DeployedAt:            big.NewInt(0),
	}

	tl1, err := NewTimeLocks(params)
	require.NoError(t, err)

	built := tl1.Build()
	tl2, err := TimeLocksFromBigInt(built)
	require.NoError(t, err)
	assert.NotNil(t, tl2)

	// Verify the round-trip by comparing the built values
	built2 := tl2.Build()
	assert.Equal(t, built.String(), built2.String())
}

func mustParseBigInt(t *testing.T, s string) *big.Int {
	val, ok := new(big.Int).SetString(s, 10)
	require.True(t, ok)
	return val
}
