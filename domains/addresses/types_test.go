package addresses

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddressLike_Interfaces(t *testing.T) {
	evmAddr := NewEvmAddress("0x0742D35CC6634c0532925A3b844bc9E7595f0Beb")
	solanaAddr, err := NewSolanaAddress("So11111111111111111111111111111111111111112")
	assert.NoError(t, err)

	// Test that both implement AddressLike
	var _ AddressLike = evmAddr
	var _ AddressLike = solanaAddr

	// Test methods exist
	assert.NotEmpty(t, evmAddr.ToString())
	assert.NotEmpty(t, solanaAddr.ToString())
	assert.NotNil(t, evmAddr.ToBuffer())
	buf := solanaAddr.ToBuffer()
	assert.NotNil(t, buf)
}
