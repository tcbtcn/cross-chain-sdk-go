package sdk

import (
	"context"

	"github.com/tcbtcn/cross-chain-sdk-go/crypto/eip712"
)

type BlockchainProvider interface {
	SignTypedData(ctx context.Context, walletAddress string, typedData eip712.TypedData) (string, error)
	EthCall(ctx context.Context, contractAddress string, callData string) (string, error)
}
