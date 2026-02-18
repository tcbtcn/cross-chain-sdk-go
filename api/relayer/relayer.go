package relayer

import (
	"context"
	"fmt"

	"github.com/tcbtcn/cross-chain-sdk-go/api/http"
)

const Version = "v1.1"

type RelayerApi struct {
	config     RelayerApiConfig
	httpClient http.Client
}

func NewRelayerApi(config RelayerApiConfig, httpClient http.Client) *RelayerApi {
	return &RelayerApi{
		config:     config,
		httpClient: httpClient,
	}
}

func (ra *RelayerApi) Submit(ctx context.Context, request interface{}) error {
	path := fmt.Sprintf("/%s/submit", Version)
	return ra.httpClient.Post(ctx, path, request, nil)
}

func (ra *RelayerApi) SubmitBatch(ctx context.Context, requests []interface{}) error {
	path := fmt.Sprintf("/%s/submit/many", Version)
	return ra.httpClient.Post(ctx, path, requests, nil)
}

func (ra *RelayerApi) SubmitSecret(ctx context.Context, orderHash, secret string) error {
	path := fmt.Sprintf("/%s/submit/secret", Version)
	request := SubmitSecretRequest{
		OrderHash: orderHash,
		Secret:    secret,
	}
	return ra.httpClient.Post(ctx, path, request, nil)
}
