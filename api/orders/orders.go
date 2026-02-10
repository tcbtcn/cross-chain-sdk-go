package orders

import (
	"context"
	"fmt"

	"github.com/dawitel/cross-chain-sdk-go/api"
	"github.com/dawitel/cross-chain-sdk-go/api/http"
	"github.com/dawitel/cross-chain-sdk-go/chains"
)

const Version = "v1.1"

type OrdersApi struct {
	config     OrdersApiConfig
	httpClient http.Client
}

func NewOrdersApi(config OrdersApiConfig, httpClient http.Client) *OrdersApi {
	return &OrdersApi{
		config:     config,
		httpClient: httpClient,
	}
}

func (oa *OrdersApi) GetActiveOrders(ctx context.Context, params ActiveOrdersRequestParams) (*ActiveOrdersResponse, error) {
	queryParams := oa.buildQueryParams(params)
	path := fmt.Sprintf("/%s/order/active/%s", Version, queryParams)

	var response ActiveOrdersResponse
	if err := oa.httpClient.Get(ctx, path, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (oa *OrdersApi) GetOrderStatus(ctx context.Context, orderHash string) (*OrderStatusResponse, error) {
	path := fmt.Sprintf("/%s/order/status/%s", Version, orderHash)

	var response OrderStatusResponse
	if err := oa.httpClient.Get(ctx, path, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (oa *OrdersApi) GetOrdersByMaker(ctx context.Context, params OrdersByMakerParams) (*OrdersByMakerResponse, error) {
	queryParams := oa.buildQueryParamsForMaker(params)
	path := fmt.Sprintf("/%s/order/maker/%s/%s", Version, params.Address, queryParams)

	var response OrdersByMakerResponse
	if err := oa.httpClient.Get(ctx, path, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (oa *OrdersApi) GetReadyToAcceptSecretFills(ctx context.Context, orderHash string) (*ReadyToAcceptSecretFills, error) {
	path := fmt.Sprintf("/%s/order/ready-to-accept-secret-fills/%s", Version, orderHash)

	var response ReadyToAcceptSecretFills
	if err := oa.httpClient.Get(ctx, path, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (oa *OrdersApi) GetReadyToExecutePublicActions(ctx context.Context) (*ReadyToExecutePublicActions, error) {
	path := fmt.Sprintf("/%s/order/ready-to-execute-public-actions", Version)

	var response ReadyToExecutePublicActions
	if err := oa.httpClient.Get(ctx, path, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (oa *OrdersApi) GetPublishedSecrets(ctx context.Context, orderHash string) (*PublishedSecretsResponse, error) {
	path := fmt.Sprintf("/%s/order/secrets/%s", Version, orderHash)

	var response PublishedSecretsResponse
	if err := oa.httpClient.Get(ctx, path, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (oa *OrdersApi) GetCancellableOrders(ctx context.Context, chainType chains.ChainType, page, limit int) (*CancellableOrdersResponse, error) {
	queryParams := map[string]interface{}{
		"chainType": string(chainType),
		"page":      page,
		"limit":     limit,
	}
	path := fmt.Sprintf("/%s/order/cancelable-by-resolvers%s", Version, api.ConcatQueryParams(queryParams))

	var response CancellableOrdersResponse
	if err := oa.httpClient.Get(ctx, path, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

func (oa *OrdersApi) buildQueryParams(params ActiveOrdersRequestParams) string {
	queryParams := map[string]interface{}{}
	if params.Page > 0 {
		queryParams["page"] = params.Page
	}
	if params.Limit > 0 {
		queryParams["limit"] = params.Limit
	}
	if params.SrcChainID != nil {
		queryParams["srcChainId"] = int(*params.SrcChainID)
	}
	if params.DstChainID != nil {
		queryParams["dstChainId"] = int(*params.DstChainID)
	}
	return api.ConcatQueryParams(queryParams)
}

func (oa *OrdersApi) buildQueryParamsForMaker(params OrdersByMakerParams) string {
	queryParams := map[string]interface{}{}
	if params.Page > 0 {
		queryParams["page"] = params.Page
	}
	if params.Limit > 0 {
		queryParams["limit"] = params.Limit
	}
	return api.ConcatQueryParams(queryParams)
}
