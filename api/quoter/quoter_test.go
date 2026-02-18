package quoter

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	clienthttp "github.com/tcbtcn/cross-chain-sdk-go/api/http"
	"github.com/tcbtcn/cross-chain-sdk-go/chains"
	"github.com/tcbtcn/cross-chain-sdk-go/domains/auction"
)

func sampleQuote() *QuoterResponse {
	return &QuoterResponse{
		QuoteID:         "test-quote-id-123",
		SrcChainID:      chains.Ethereum,
		DstChainID:      chains.Polygon,
		SrcTokenAddress: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48",
		DstTokenAddress: "0x2791bca1f2de4661ed88a30c99a7a9449aa84174",
		SrcTokenAmount:  "100000000",
		DstTokenAmount:  "99500000",
		Presets: QuoterPresets{
			Fast: PresetData{
				AuctionDuration:    300,
				StartAuctionIn:     0,
				InitialRateBump:    100,
				AuctionStartAmount: "100000000",
				StartAmount:        "100000000",
				AuctionEndAmount:   "99500000",
				CostInDstToken:     "500000",
				Points:             []auction.AuctionPoint{},
				AllowPartialFills:  true,
				AllowMultipleFills: true,
				GasCost: GasCost{
					GasBumpEstimate:  100000,
					GasPriceEstimate: "20000000000",
				},
				ExclusiveResolver: "0x0000000000000000000000000000000000000000",
				SecretsCount:      1,
			},
			Medium: PresetData{
				AuctionDuration:    600,
				SecretsCount:       1,
				AllowMultipleFills: true,
				AllowPartialFills:  true,
			},
			Slow: PresetData{
				AuctionDuration:    900,
				SecretsCount:       1,
				AllowMultipleFills: true,
				AllowPartialFills:  true,
			},
		},
		RecommendedPreset: PresetFast,
		SrcEscrowFactory:  "0x1111111111111111111111111111111111111111",
		DstEscrowFactory:  "0x2222222222222222222222222222222222222222",
		TimeLocks: TimeLocksRaw{
			SrcWithdrawal:         3600,
			SrcPublicWithdrawal:   7200,
			SrcCancellation:       1800,
			SrcPublicCancellation: 3600,
			DstWithdrawal:         3600,
			DstPublicWithdrawal:   7200,
			DstCancellation:       1800,
		},
		SrcSafetyDeposit: "1000000",
		DstSafetyDeposit: "1000000",
		AutoK:            1,
	}
}

func TestNewQuoterApi(t *testing.T) {
	httpClient := clienthttp.NewHTTPClient("https://api.example.com", "test-key")
	api := NewQuoterApi(QuoterApiConfig{
		URL:     "https://api.example.com",
		AuthKey: "test-key",
	}, httpClient)

	assert.NotNil(t, api)
}

func TestQuoterApi_GetQuote(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		_ = json.NewEncoder(w).Encode(sampleQuote())
	}))
	defer server.Close()

	httpClient := clienthttp.NewHTTPClient(server.URL, "")
	api := NewQuoterApi(QuoterApiConfig{
		URL:     server.URL,
		AuthKey: "",
	}, httpClient)

	params := QuoterRequestParams{
		SrcChain:        chains.Ethereum,
		DstChain:        chains.Polygon,
		SrcTokenAddress: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48",
		DstTokenAddress: "0x2791bca1f2de4661ed88a30c99a7a9449aa84174",
		Amount:          "100000000",
		WalletAddress:   "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb",
		EnableEstimate:  true,
	}

	quote, err := api.GetQuote(context.Background(), params)
	require.NoError(t, err)
	assert.NotNil(t, quote)
	assert.Equal(t, params.SrcChain, quote.SrcChainID)
	assert.Equal(t, params.DstChain, quote.DstChainID)
	assert.Equal(t, params.SrcTokenAddress, quote.SrcTokenAddress)
	assert.Equal(t, params.DstTokenAddress, quote.DstTokenAddress)
}

func TestQuoterApi_GetQuoteWithCustomPreset(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		_ = json.NewEncoder(w).Encode(sampleQuote())
	}))
	defer server.Close()

	httpClient := clienthttp.NewHTTPClient(server.URL, "")
	api := NewQuoterApi(QuoterApiConfig{
		URL:     server.URL,
		AuthKey: "",
	}, httpClient)

	params := QuoterRequestParams{
		SrcChain:        chains.Ethereum,
		DstChain:        chains.Polygon,
		SrcTokenAddress: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48",
		DstTokenAddress: "0x2791bca1f2de4661ed88a30c99a7a9449aa84174",
		Amount:          "100000000",
		WalletAddress:   "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb",
		EnableEstimate:  true,
	}

	customPreset := CustomPreset{
		AuctionDuration:    300,
		AuctionStartAmount: "100000000",
		AuctionEndAmount:   "99500000",
	}

	quote, err := api.GetQuoteWithCustomPreset(context.Background(), params, customPreset)
	require.NoError(t, err)
	assert.NotNil(t, quote)
}

func TestQuoterApi_GetQuoteWithCustomPreset_Invalid(t *testing.T) {
	httpClient := clienthttp.NewHTTPClient("https://api.example.com", "")
	api := NewQuoterApi(QuoterApiConfig{
		URL:     "https://api.example.com",
		AuthKey: "",
	}, httpClient)

	params := QuoterRequestParams{
		SrcChain: chains.Ethereum,
		DstChain: chains.Polygon,
	}

	tests := []struct {
		name         string
		customPreset CustomPreset
		wantErr      bool
	}{
		{
			name: "invalid zero duration",
			customPreset: CustomPreset{
				AuctionDuration:    0,
				AuctionStartAmount: "100000000",
				AuctionEndAmount:   "99500000",
			},
			wantErr: true,
		},
		{
			name: "invalid negative duration",
			customPreset: CustomPreset{
				AuctionDuration:    -1,
				AuctionStartAmount: "100000000",
				AuctionEndAmount:   "99500000",
			},
			wantErr: true,
		},
		{
			name: "missing start amount",
			customPreset: CustomPreset{
				AuctionDuration:  300,
				AuctionEndAmount: "99500000",
			},
			wantErr: true,
		},
		{
			name: "missing end amount",
			customPreset: CustomPreset{
				AuctionDuration:    300,
				AuctionStartAmount: "100000000",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			quote, err := api.GetQuoteWithCustomPreset(context.Background(), params, tt.customPreset)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, quote)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestQuoterApi_buildQueryParams(t *testing.T) {
	api := NewQuoterApi(QuoterApiConfig{}, clienthttp.NewHTTPClient("", ""))

	params := QuoterRequestParams{
		SrcChain:        chains.Ethereum,
		DstChain:        chains.Polygon,
		SrcTokenAddress: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48",
		DstTokenAddress: "0x2791bca1f2de4661ed88a30c99a7a9449aa84174",
		Amount:          "100000000",
		WalletAddress:   "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb",
		EnableEstimate:  true,
		Permit:          "0x123",
		Fee:             100,
		Source:          "test",
		IsPermit2:       true,
	}

	queryParams := api.buildQueryParams(params)
	assert.NotEmpty(t, queryParams)
	assert.Contains(t, queryParams, "srcChainId")
	assert.Contains(t, queryParams, "dstChainId")
	assert.Contains(t, queryParams, "enableEstimate")
}
