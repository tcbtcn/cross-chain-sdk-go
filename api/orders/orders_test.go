package orders

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	clienthttp "github.com/dawitel/cross-chain-sdk-go/api/http"
	"github.com/dawitel/cross-chain-sdk-go/chains"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewOrdersApi(t *testing.T) {
	httpClient := clienthttp.NewHTTPClient("https://api.example.com", "test-key")
	api := NewOrdersApi(OrdersApiConfig{
		URL:     "https://api.example.com",
		AuthKey: "test-key",
	}, httpClient)

	assert.NotNil(t, api)
}

func TestOrdersApi_GetOrderStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		_ = json.NewEncoder(w).Encode(OrderStatusResponse{
			OrderHash:  "0x123",
			Status:     OrderStatusPending,
			SrcChainID: chains.Ethereum,
			DstChainID: chains.Polygon,
		})
	}))
	defer server.Close()

	httpClient := clienthttp.NewHTTPClient(server.URL, "")
	api := NewOrdersApi(OrdersApiConfig{
		URL:     server.URL,
		AuthKey: "",
	}, httpClient)

	status, err := api.GetOrderStatus(context.Background(), "0x123")
	require.NoError(t, err)
	assert.NotNil(t, status)
	assert.Equal(t, "0x123", status.OrderHash)
}

func TestOrdersApi_GetActiveOrders(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		_ = json.NewEncoder(w).Encode(ActiveOrdersResponse{
			Items:      []ActiveOrder{},
			TotalCount: 0,
		})
	}))
	defer server.Close()

	httpClient := clienthttp.NewHTTPClient(server.URL, "")
	api := NewOrdersApi(OrdersApiConfig{
		URL:     server.URL,
		AuthKey: "",
	}, httpClient)

	params := ActiveOrdersRequestParams{
		Page:  1,
		Limit: 10,
	}

	orders, err := api.GetActiveOrders(context.Background(), params)
	require.NoError(t, err)
	assert.NotNil(t, orders)
}

func TestOrdersApi_GetOrdersByMaker(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		_ = json.NewEncoder(w).Encode(OrdersByMakerResponse{
			Items:      []OrderFillsByMakerOutput{},
			TotalCount: 0,
		})
	}))
	defer server.Close()

	httpClient := clienthttp.NewHTTPClient(server.URL, "")
	api := NewOrdersApi(OrdersApiConfig{
		URL:     server.URL,
		AuthKey: "",
	}, httpClient)

	params := OrdersByMakerParams{
		Address: "0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb",
		Page:    1,
		Limit:   10,
	}

	orders, err := api.GetOrdersByMaker(context.Background(), params)
	require.NoError(t, err)
	assert.NotNil(t, orders)
}

func TestOrdersApi_GetReadyToAcceptSecretFills(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		_ = json.NewEncoder(w).Encode(ReadyToAcceptSecretFills{
			Fills: []ReadyToAcceptSecretFill{},
		})
	}))
	defer server.Close()

	httpClient := clienthttp.NewHTTPClient(server.URL, "")
	api := NewOrdersApi(OrdersApiConfig{
		URL:     server.URL,
		AuthKey: "",
	}, httpClient)

	fills, err := api.GetReadyToAcceptSecretFills(context.Background(), "0x123")
	require.NoError(t, err)
	assert.NotNil(t, fills)
}

func TestOrdersApi_GetPublishedSecrets(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		_ = json.NewEncoder(w).Encode(PublishedSecretsResponse{
			Secrets: []PublicSecret{},
		})
	}))
	defer server.Close()

	httpClient := clienthttp.NewHTTPClient(server.URL, "")
	api := NewOrdersApi(OrdersApiConfig{
		URL:     server.URL,
		AuthKey: "",
	}, httpClient)

	secrets, err := api.GetPublishedSecrets(context.Background(), "0x123")
	require.NoError(t, err)
	assert.NotNil(t, secrets)
}

func TestOrdersApi_GetCancellableOrders(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		_ = json.NewEncoder(w).Encode(CancellableOrdersResponse{
			Items:      []interface{}{},
			TotalCount: 0,
		})
	}))
	defer server.Close()

	httpClient := clienthttp.NewHTTPClient(server.URL, "")
	api := NewOrdersApi(OrdersApiConfig{
		URL:     server.URL,
		AuthKey: "",
	}, httpClient)

	orders, err := api.GetCancellableOrders(context.Background(), chains.ChainTypeEVM, 1, 10)
	require.NoError(t, err)
	assert.NotNil(t, orders)
}

func TestOrdersApi_GetReadyToExecutePublicActions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		_ = json.NewEncoder(w).Encode(ReadyToExecutePublicActions{
			Actions: []ReadyToExecutePublicAction{},
		})
	}))
	defer server.Close()

	httpClient := clienthttp.NewHTTPClient(server.URL, "")
	api := NewOrdersApi(OrdersApiConfig{
		URL:     server.URL,
		AuthKey: "",
	}, httpClient)

	actions, err := api.GetReadyToExecutePublicActions(context.Background())
	require.NoError(t, err)
	assert.NotNil(t, actions)
}
