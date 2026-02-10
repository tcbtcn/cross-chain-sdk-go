package quoter

import (
	"context"
	"fmt"

	"github.com/dawitel/cross-chain-sdk-go/api"
	"github.com/dawitel/cross-chain-sdk-go/api/http"
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
	return nil
}
