package quoter

import (
	"context"
	"fmt"
	"math/big"

	"github.com/dawitel/cross-chain-sdk-go/api"
	"github.com/dawitel/cross-chain-sdk-go/api/http"
	"github.com/dawitel/cross-chain-sdk-go/chains"
	"github.com/dawitel/cross-chain-sdk-go/domains/addresses"
)

const Version = "v1.1"

type QuoterApi struct {
	config     QuoterApiConfig
	httpClient http.Client
}

func NewQuoterApi(config QuoterApiConfig, httpClient http.Client) *QuoterApi {
	return &QuoterApi{
		config:     config,
		httpClient: httpClient,
	}
}

func (qa *QuoterApi) GetQuote(ctx context.Context, params QuoterRequestParams) (*QuoterResponse, error) {
	if err := qa.validateQuoteParams(params); err != nil {
		return nil, err
	}

	queryParams := qa.buildQueryParams(params)
	path := fmt.Sprintf("/%s/quote/receive/%s", Version, queryParams)

	var response QuoterResponse
	if err := qa.httpClient.Get(ctx, path, &response); err != nil {
		return nil, err
	}

	response.SrcChainID = params.SrcChain
	response.DstChainID = params.DstChain
	response.SrcTokenAddress = params.SrcTokenAddress
	response.DstTokenAddress = params.DstTokenAddress

	return &response, nil
}

func (qa *QuoterApi) GetQuoteWithCustomPreset(ctx context.Context, params QuoterRequestParams, customPreset CustomPreset) (*QuoterResponse, error) {
	if err := qa.validateQuoteParams(params); err != nil {
		return nil, err
	}
	if err := qa.validateCustomPreset(customPreset); err != nil {
		return nil, err
	}

	queryParams := qa.buildQueryParams(params)
	path := fmt.Sprintf("/%s/quote/receive/%s", Version, queryParams)

	var response QuoterResponse
	if err := qa.httpClient.Post(ctx, path, customPreset, &response); err != nil {
		return nil, err
	}

	response.SrcChainID = params.SrcChain
	response.DstChainID = params.DstChain
	response.SrcTokenAddress = params.SrcTokenAddress
	response.DstTokenAddress = params.DstTokenAddress

	return &response, nil
}

func (qa *QuoterApi) buildQueryParams(params QuoterRequestParams) string {
	queryParams := map[string]interface{}{
		"srcChainId":      int(params.SrcChain),
		"dstChainId":      int(params.DstChain),
		"srcTokenAddress": params.SrcTokenAddress,
		"dstTokenAddress": params.DstTokenAddress,
		"amount":          params.Amount,
		"walletAddress":   params.WalletAddress,
	}

	if params.EnableEstimate {
		queryParams["enableEstimate"] = true
	}
	if params.Permit != "" {
		queryParams["permit"] = params.Permit
	}
	if params.Fee > 0 {
		queryParams["fee"] = params.Fee
	}
	if params.Source != "" {
		queryParams["source"] = params.Source
	}
	if params.IsPermit2 {
		queryParams["isPermit2"] = true
	}

	return api.ConcatQueryParams(queryParams)
}

func (qa *QuoterApi) validateQuoteParams(params QuoterRequestParams) error {
	if params.SrcChain == params.DstChain {
		return fmt.Errorf("srcChain and dstChain should be different")
	}

	if params.Amount == "" {
		return fmt.Errorf("amount is required")
	}
	amount, ok := new(big.Int).SetString(params.Amount, 10)
	if !ok || amount.Sign() <= 0 {
		return fmt.Errorf("%s is invalid amount", params.Amount)
	}

	if chains.IsEvm(params.SrcChain) {
		if params.WalletAddress != "" {
			_ = addresses.NewEvmAddress(params.WalletAddress)
		}

		if chains.IsEvm(params.DstChain) {
			if params.DstTokenAddress == addresses.ZeroAddress {
				return fmt.Errorf("replace %s with %s", addresses.ZeroAddress, addresses.NativeCurrency)
			}
		}
	}

	if params.Fee > 0 && params.Source == "" {
		return fmt.Errorf("cannot use fee without source")
	}

	return nil
}

func (qa *QuoterApi) validateCustomPreset(customPreset CustomPreset) error {
	if customPreset.AuctionDuration <= 0 {
		return fmt.Errorf("auctionDuration must be positive")
	}
	if customPreset.AuctionStartAmount == "" {
		return fmt.Errorf("auctionStartAmount is required")
	}
	if customPreset.AuctionEndAmount == "" {
		return fmt.Errorf("auctionEndAmount is required")
	}

	startAmount, ok := new(big.Int).SetString(customPreset.AuctionStartAmount, 10)
	if !ok || startAmount.Sign() < 0 {
		return fmt.Errorf("Invalid auctionStartAmount")
	}

	endAmount, ok := new(big.Int).SetString(customPreset.AuctionEndAmount, 10)
	if !ok || endAmount.Sign() < 0 {
		return fmt.Errorf("Invalid auctionEndAmount")
	}

	if len(customPreset.Points) > 0 {
		for i, point := range customPreset.Points {
			pointAmount, ok := new(big.Int).SetString(point.ToTokenAmount, 10)
			if !ok || pointAmount.Sign() < 0 {
				return fmt.Errorf("Invalid toTokenAmount in point %d", i)
			}
			if pointAmount.Cmp(startAmount) < 0 || pointAmount.Cmp(endAmount) > 0 {
				return fmt.Errorf("Point %d toTokenAmount must be between auctionStartAmount and auctionEndAmount", i)
			}
			if point.Delay < 0 {
				return fmt.Errorf("Point %d delay must be non-negative", i)
			}
		}
	}

	return nil
}
