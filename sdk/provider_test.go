package sdk

import (
	"context"
	"testing"

	"github.com/dawitel/cross-chain-sdk-go/crypto/eip712"
	"github.com/dawitel/cross-chain-sdk-go/testutils"
	"github.com/stretchr/testify/assert"
)

func TestBlockchainProvider_Interface(t *testing.T) {
	mockProvider := &testutils.MockBlockchainProvider{}

	var _ BlockchainProvider = mockProvider

	// Test interface methods exist
	_, err := mockProvider.SignTypedData(context.Background(), "0x123", eip712.TypedData{})
	assert.NoError(t, err) // Mock returns nil error by default

	_, err = mockProvider.EthCall(context.Background(), "0x123", "0x456")
	assert.NoError(t, err)
}
