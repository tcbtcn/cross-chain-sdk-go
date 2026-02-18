package quoter

import (
	"github.com/tcbtcn/cross-chain-sdk-go/chains"
	"github.com/tcbtcn/cross-chain-sdk-go/domains/auction"
)

type QuoterRequestParams struct {
	SrcChain        chains.SupportedChain
	DstChain        chains.SupportedChain
	SrcTokenAddress string
	DstTokenAddress string
	Amount          string
	WalletAddress   string
	EnableEstimate  bool
	Permit          string
	Fee             int
	Source          string
	IsPermit2       bool
}

type QuoterResponse struct {
	QuoteID            string                `json:"quoteId"`
	SrcChainID         chains.SupportedChain `json:"srcChainId"`
	DstChainID         chains.SupportedChain `json:"dstChainId"`
	SrcTokenAddress    string                `json:"srcTokenAddress"`
	DstTokenAddress    string                `json:"dstTokenAddress"`
	SrcTokenAmount     string                `json:"srcTokenAmount"`
	DstTokenAmount     string                `json:"dstTokenAmount"`
	Presets            QuoterPresets         `json:"presets"`
	SrcEscrowFactory   string                `json:"srcEscrowFactory"`
	DstEscrowFactory   string                `json:"dstEscrowFactory"`
	RecommendedPreset  PresetEnum            `json:"recommendedPreset"`
	Prices             Cost                  `json:"prices"`
	Volume             Cost                  `json:"volume"`
	Whitelist          []string              `json:"whitelist"`
	TimeLocks          TimeLocksRaw          `json:"timeLocks"`
	SrcSafetyDeposit   string                `json:"srcSafetyDeposit"`
	DstSafetyDeposit   string                `json:"dstSafetyDeposit"`
	AutoK              int                   `json:"autoK"`
	NativeOrderFactory string                `json:"nativeOrderFactoryAddress,omitempty"`
	NativeOrderImpl    string                `json:"nativeOrderImplAddress,omitempty"`
}

type QuoterPresets struct {
	Fast   PresetData  `json:"fast"`
	Medium PresetData  `json:"medium"`
	Slow   PresetData  `json:"slow"`
	Custom *PresetData `json:"custom,omitempty"`
}

type PresetData struct {
	AuctionDuration    int                    `json:"auctionDuration"`
	StartAuctionIn     int                    `json:"startAuctionIn"`
	InitialRateBump    int                    `json:"initialRateBump"`
	AuctionStartAmount string                 `json:"auctionStartAmount"`
	StartAmount        string                 `json:"startAmount"`
	AuctionEndAmount   string                 `json:"auctionEndAmount"`
	CostInDstToken     string                 `json:"costInDstToken"`
	Points             []auction.AuctionPoint `json:"points"`
	AllowPartialFills  bool                   `json:"allowPartialFills"`
	AllowMultipleFills bool                   `json:"allowMultipleFills"`
	GasCost            GasCost                `json:"gasCost"`
	ExclusiveResolver  string                 `json:"exclusiveResolver"`
	SecretsCount       int                    `json:"secretsCount"`
}

type GasCost struct {
	GasBumpEstimate  int    `json:"gasBumpEstimate"`
	GasPriceEstimate string `json:"gasPriceEstimate"`
}

type Cost struct {
	USD USDCost `json:"usd"`
}

type USDCost struct {
	SrcToken string `json:"srcToken"`
	DstToken string `json:"dstToken"`
}

type TimeLocksRaw struct {
	SrcWithdrawal         int `json:"srcWithdrawal"`
	SrcPublicWithdrawal   int `json:"srcPublicWithdrawal"`
	SrcCancellation       int `json:"srcCancellation"`
	SrcPublicCancellation int `json:"srcPublicCancellation"`
	DstWithdrawal         int `json:"dstWithdrawal"`
	DstPublicWithdrawal   int `json:"dstPublicWithdrawal"`
	DstCancellation       int `json:"dstCancellation"`
}

type PresetEnum string

const (
	PresetFast   PresetEnum = "fast"
	PresetMedium PresetEnum = "medium"
	PresetSlow   PresetEnum = "slow"
	PresetCustom PresetEnum = "custom"
)

type CustomPreset struct {
	AuctionDuration    int                 `json:"auctionDuration"`
	AuctionStartAmount string              `json:"auctionStartAmount"`
	AuctionEndAmount   string              `json:"auctionEndAmount"`
	Points             []CustomPresetPoint `json:"points,omitempty"`
}

type CustomPresetPoint struct {
	ToTokenAmount string `json:"toTokenAmount"`
	Delay         int    `json:"delay"`
}

type QuoterApiConfig struct {
	URL     string
	AuthKey string
}
