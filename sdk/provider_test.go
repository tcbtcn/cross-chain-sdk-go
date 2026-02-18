package sdk

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tcbtcn/cross-chain-sdk-go/crypto/eip712"
	"github.com/tcbtcn/cross-chain-sdk-go/testutils"
)

func TestBlockchainProvider_Interface(t *testing.T) {
	mockProvider := &testutils.MockBlockchainProvider{}

	var _ BlockchainProvider = mockProvider

	mockProvider.On("SignTypedData", context.Background(), "0x123", eip712.TypedData{}).Return("0xsignature", nil)
	_, err := mockProvider.SignTypedData(context.Background(), "0x123", eip712.TypedData{})
	assert.NoError(t, err)

	mockProvider.On("EthCall", context.Background(), "0x123", "0x456").Return("0xresult", nil)
	_, err = mockProvider.EthCall(context.Background(), "0x123", "0x456")
	assert.NoError(t, err)
}
