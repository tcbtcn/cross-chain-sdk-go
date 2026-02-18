package testutils

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/tcbtcn/cross-chain-sdk-go/api/http"
	"github.com/tcbtcn/cross-chain-sdk-go/crypto/eip712"
)

// MockHTTPClient is a mock implementation of http.Client
type MockHTTPClient struct {
	mock.Mock
}

func (m *MockHTTPClient) Get(ctx context.Context, url string, result interface{}) error {
	args := m.Called(ctx, url, result)
	_ = args.Get(0)
	_ = result
	return args.Error(1)
}

func (m *MockHTTPClient) Post(ctx context.Context, url string, data interface{}, result interface{}) error {
	args := m.Called(ctx, url, data, result)
	_ = args.Get(0)
	_ = result
	return args.Error(1)
}

// MockBlockchainProvider is a mock implementation of sdk.BlockchainProvider
type MockBlockchainProvider struct {
	mock.Mock
}

func (m *MockBlockchainProvider) SignTypedData(ctx context.Context, address string, data eip712.TypedData) (string, error) {
	args := m.Called(ctx, address, data)
	return args.String(0), args.Error(1)
}

func (m *MockBlockchainProvider) EthCall(ctx context.Context, contractAddress string, callData string) (string, error) {
	args := m.Called(ctx, contractAddress, callData)
	return args.String(0), args.Error(1)
}

// Ensure interfaces are satisfied
var _ http.Client = (*MockHTTPClient)(nil)
