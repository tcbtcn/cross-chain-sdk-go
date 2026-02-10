package http

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewHTTPClient(t *testing.T) {
	client := NewHTTPClient("https://api.example.com", "test-key")
	assert.NotNil(t, client)
	assert.Equal(t, "https://api.example.com", client.baseURL)
	assert.Equal(t, "test-key", client.authKey)
	assert.NotNil(t, client.client)
}

func TestHTTPClient_Get(t *testing.T) {
	tests := []struct {
		name         string
		statusCode   int
		responseBody interface{}
		authKey      string
		wantErr      bool
		checkAuth    bool
	}{
		{
			name:       "successful GET",
			statusCode: 200,
			responseBody: map[string]string{
				"result": "success",
			},
			authKey:   "test-key",
			wantErr:   false,
			checkAuth: true,
		},
		{
			name:       "GET without auth",
			statusCode: 200,
			responseBody: map[string]string{
				"result": "success",
			},
			authKey:   "",
			wantErr:   false,
			checkAuth: false,
		},
		{
			name:       "GET with 404",
			statusCode: 404,
			responseBody: map[string]string{
				"error": "not found",
			},
			authKey:   "test-key",
			wantErr:   true,
			checkAuth: true,
		},
		{
			name:       "GET with 500",
			statusCode: 500,
			responseBody: map[string]string{
				"error": "internal server error",
			},
			authKey:   "test-key",
			wantErr:   true,
			checkAuth: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.checkAuth {
					auth := r.Header.Get("Authorization")
					assert.Equal(t, "Bearer "+tt.authKey, auth)
				}
				assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
				assert.Equal(t, "GET", r.Method)

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				_ = json.NewEncoder(w).Encode(tt.responseBody)
			}))
			defer server.Close()

			client := NewHTTPClient(server.URL, tt.authKey)
			var result map[string]string
			err := client.Get(context.Background(), "", &result)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.responseBody.(map[string]string)["result"], result["result"])
			}
		})
	}
}

func TestHTTPClient_Post(t *testing.T) {
	tests := []struct {
		name         string
		statusCode   int
		requestData  interface{}
		responseBody interface{}
		authKey      string
		wantErr      bool
		checkAuth    bool
	}{
		{
			name:       "successful POST",
			statusCode: 200,
			requestData: map[string]string{
				"data": "test",
			},
			responseBody: map[string]string{
				"result": "success",
			},
			authKey:   "test-key",
			wantErr:   false,
			checkAuth: true,
		},
		{
			name:       "POST without auth",
			statusCode: 200,
			requestData: map[string]string{
				"data": "test",
			},
			responseBody: map[string]string{
				"result": "success",
			},
			authKey:   "",
			wantErr:   false,
			checkAuth: false,
		},
		{
			name:       "POST with 400",
			statusCode: 400,
			requestData: map[string]string{
				"data": "invalid",
			},
			responseBody: map[string]string{
				"error": "bad request",
			},
			authKey:   "test-key",
			wantErr:   true,
			checkAuth: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.checkAuth {
					auth := r.Header.Get("Authorization")
					assert.Equal(t, "Bearer "+tt.authKey, auth)
				}
				assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
				assert.Equal(t, "POST", r.Method)

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				_ = json.NewEncoder(w).Encode(tt.responseBody)
			}))
			defer server.Close()

			client := NewHTTPClient(server.URL, tt.authKey)
			var result map[string]string
			err := client.Post(context.Background(), "", tt.requestData, &result)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.responseBody.(map[string]string)["result"], result["result"])
			}
		})
	}
}

func TestHTTPClient_ContextCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(100 * time.Millisecond)
		w.WriteHeader(200)
	}))
	defer server.Close()

	client := NewHTTPClient(server.URL, "test-key")
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	var result map[string]string
	err := client.Get(ctx, "", &result)
	assert.Error(t, err)
}

func TestHTTPClient_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		_, _ = w.Write([]byte("invalid json"))
	}))
	defer server.Close()

	client := NewHTTPClient(server.URL, "test-key")
	var result map[string]string
	err := client.Get(context.Background(), "", &result)
	assert.Error(t, err)
}

func TestHTTPClient_EmptyResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
	}))
	defer server.Close()

	client := NewHTTPClient(server.URL, "test-key")
	err := client.Get(context.Background(), "", nil)
	assert.NoError(t, err)
}
