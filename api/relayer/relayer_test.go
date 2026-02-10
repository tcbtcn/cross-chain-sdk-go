package relayer

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	clienthttp "github.com/dawitel/cross-chain-sdk-go/api/http"
	"github.com/stretchr/testify/assert"
)

func TestNewRelayerApi(t *testing.T) {
	httpClient := clienthttp.NewHTTPClient("https://api.example.com", "test-key")
	api := NewRelayerApi(RelayerApiConfig{
		URL:     "https://api.example.com",
		AuthKey: "test-key",
	}, httpClient)

	assert.NotNil(t, api)
}

func TestRelayerApi_Submit(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
	}))
	defer server.Close()

	httpClient := clienthttp.NewHTTPClient(server.URL, "")
	api := NewRelayerApi(RelayerApiConfig{
		URL:     server.URL,
		AuthKey: "",
	}, httpClient)

	request := RelayerRequestEvm{
		SrcChainID: 1,
		Order:      map[string]interface{}{},
		Signature:  "0x123",
		QuoteID:    "quote-123",
	}

	err := api.Submit(context.Background(), request)
	assert.NoError(t, err)
}

func TestRelayerApi_SubmitBatch(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
	}))
	defer server.Close()

	httpClient := clienthttp.NewHTTPClient(server.URL, "")
	api := NewRelayerApi(RelayerApiConfig{
		URL:     server.URL,
		AuthKey: "",
	}, httpClient)

	requests := []interface{}{
		RelayerRequestEvm{SrcChainID: 1},
		RelayerRequestEvm{SrcChainID: 137},
	}

	err := api.SubmitBatch(context.Background(), requests)
	assert.NoError(t, err)
}

func TestRelayerApi_SubmitSecret(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)

		var req SubmitSecretRequest
		json.NewDecoder(r.Body).Decode(&req)
		assert.Equal(t, "0x123", req.OrderHash)
		assert.Equal(t, "0xsecret", req.Secret)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
	}))
	defer server.Close()

	httpClient := clienthttp.NewHTTPClient(server.URL, "")
	api := NewRelayerApi(RelayerApiConfig{
		URL:     server.URL,
		AuthKey: "",
	}, httpClient)

	err := api.SubmitSecret(context.Background(), "0x123", "0xsecret")
	assert.NoError(t, err)
}

func TestRelayerApi_Submit_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(map[string]string{"error": "bad request"})
	}))
	defer server.Close()

	httpClient := clienthttp.NewHTTPClient(server.URL, "")
	api := NewRelayerApi(RelayerApiConfig{
		URL:     server.URL,
		AuthKey: "",
	}, httpClient)

	request := RelayerRequestEvm{
		SrcChainID: 1,
	}

	err := api.Submit(context.Background(), request)
	assert.Error(t, err)
}
