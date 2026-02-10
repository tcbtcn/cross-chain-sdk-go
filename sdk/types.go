package sdk

import (
	"math/big"

	"github.com/dawitel/cross-chain-sdk-go/chains"
	"github.com/dawitel/cross-chain-sdk-go/domains/addresses"
	"github.com/dawitel/cross-chain-sdk-go/domains/hashlock"
	"github.com/dawitel/cross-chain-sdk-go/orders"
)

type Config struct {
	URL                string
	AuthKey            string
	BlockchainProvider BlockchainProvider
	HTTPClient         interface{}
}

type QuoteParams struct {
	SrcChainID      chains.SupportedChain
	DstChainID      chains.SupportedChain
	SrcTokenAddress string
	DstTokenAddress string
	Amount          string
	WalletAddress   string
	EnableEstimate  bool
	Permit          string
	TakingFeeBps    int
	Source          string
	IsPermit2       bool
}

type OrderParams struct {
	WalletAddress string
	HashLock      *hashlock.HashLock
	SecretHashes  []string
	Permit        string
	Receiver      string
	Preset        string
	Nonce         *big.Int
	Fee           *TakingFeeInfo
	Source        string
	IsPermit2     bool
	CustomPreset  interface{}
}

type TakingFeeInfo struct {
	TakingFeeBps      int
	TakingFeeReceiver string
}

type OrderInfo struct {
	Order     orders.LimitOrderV4Struct
	Signature string
	QuoteID   string
	OrderHash string
	Extension string
}

type PreparedOrder struct {
	Order   interface{}
	Hash    string
	QuoteID string
}

type QuoteCustomPresetParams struct {
	CustomPreset interface{}
}

type EvmOrderCancellationData struct {
	OrderHash            string
	Maker                *addresses.EvmAddress
	SrcChainID           chains.SupportedChain
	DstChainID           chains.SupportedChain
	Order                orders.LimitOrderV4Struct
	Extension            string
	RemainingMakerAmount *big.Int
}

type SvmOrderCancellationData struct {
	OrderHash          []byte
	Maker              *addresses.SolanaAddress
	Token              *addresses.SolanaAddress
	CancellationConfig interface{} // ResolverCancellationConfig - simplified for now
	IsAssetNative      bool
}
