package wsapi

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNewWebSocketApi_WithProvider(t *testing.T) {
	mockProvider := new(MockWsProviderConnector)

	ws, err := NewWebSocketApi(mockProvider)
	require.NoError(t, err)
	assert.NotNil(t, ws)
	assert.Equal(t, mockProvider, ws.provider)
	assert.NotNil(t, ws.RPC)
	assert.NotNil(t, ws.Order)
}

func TestNewWebSocketApi_WithConfig(t *testing.T) {
	config := WsApiConfig{
		URL:      "wss://api.1inch.dev/fusion-plus/ws",
		AuthKey:  "test-key",
		LazyInit: true,
	}

	ws, err := NewWebSocketApi(config)
	require.NoError(t, err)
	assert.NotNil(t, ws)
	assert.NotNil(t, ws.RPC)
	assert.NotNil(t, ws.Order)
}

func TestNewWebSocketApi_WithConfig_HTTPURL(t *testing.T) {
	config := WsApiConfig{
		URL:      "https://api.1inch.dev/fusion-plus/ws",
		AuthKey:  "test-key",
		LazyInit: true,
	}

	ws, err := NewWebSocketApi(config)
	require.NoError(t, err)
	assert.NotNil(t, ws)
}

func TestNewWebSocketApi_WithConfig_HTTPURL_NoProtocol(t *testing.T) {
	config := WsApiConfig{
		URL:      "api.1inch.dev/fusion-plus/ws",
		AuthKey:  "test-key",
		LazyInit: true,
	}

	ws, err := NewWebSocketApi(config)
	require.NoError(t, err)
	assert.NotNil(t, ws)
}

func TestNewWebSocketApi_InvalidType(t *testing.T) {
	ws, err := NewWebSocketApi("invalid")
	assert.Error(t, err)
	assert.Nil(t, ws)
	assert.Contains(t, err.Error(), "invalid config or provider type")
}

func TestWebSocketApi_Init(t *testing.T) {
	mockProvider := new(MockWsProviderConnector)
	mockProvider.Mock.On("Init")

	ws, err := NewWebSocketApi(mockProvider)
	require.NoError(t, err)

	ws.Init()
	mockProvider.AssertCalled(t, "Init")
}

func TestWebSocketApi_On(t *testing.T) {
	mockProvider := new(MockWsProviderConnector)
	cb := func() {}
	mockProvider.Mock.On("On", EventOpen, mock.AnythingOfType("func()"))

	ws, err := NewWebSocketApi(mockProvider)
	require.NoError(t, err)

	ws.On(EventOpen, cb)
	mockProvider.AssertCalled(t, "On", EventOpen, mock.AnythingOfType("func()"))
}

func TestWebSocketApi_Off(t *testing.T) {
	mockProvider := new(MockWsProviderConnector)
	cb := func() {}
	mockProvider.Mock.On("Off", EventOpen, mock.AnythingOfType("func()"))

	ws, err := NewWebSocketApi(mockProvider)
	require.NoError(t, err)

	ws.Off(EventOpen, cb)
	mockProvider.AssertCalled(t, "Off", EventOpen, mock.AnythingOfType("func()"))
}

func TestWebSocketApi_OnOpen(t *testing.T) {
	mockProvider := new(MockWsProviderConnector)
	cb := func() {}
	mockProvider.Mock.On("OnOpen", mock.AnythingOfType("wsapi.OnOpenCb"))

	ws, err := NewWebSocketApi(mockProvider)
	require.NoError(t, err)

	ws.OnOpen(cb)
	mockProvider.AssertCalled(t, "OnOpen", mock.AnythingOfType("wsapi.OnOpenCb"))
}

func TestWebSocketApi_OnClose(t *testing.T) {
	mockProvider := new(MockWsProviderConnector)
	cb := func() {}
	mockProvider.Mock.On("OnClose", mock.AnythingOfType("wsapi.OnCloseCb"))

	ws, err := NewWebSocketApi(mockProvider)
	require.NoError(t, err)

	ws.OnClose(cb)
	mockProvider.AssertCalled(t, "OnClose", mock.AnythingOfType("wsapi.OnCloseCb"))
}

func TestWebSocketApi_OnError(t *testing.T) {
	mockProvider := new(MockWsProviderConnector)
	cb := func(err error) {}
	mockProvider.Mock.On("OnError", mock.AnythingOfType("wsapi.OnErrorCb"))

	ws, err := NewWebSocketApi(mockProvider)
	require.NoError(t, err)

	ws.OnError(cb)
	mockProvider.AssertCalled(t, "OnError", mock.AnythingOfType("wsapi.OnErrorCb"))
}

func TestWebSocketApi_OnMessage(t *testing.T) {
	mockProvider := new(MockWsProviderConnector)
	cb := func(data interface{}) {}
	mockProvider.Mock.On("OnMessage", mock.AnythingOfType("wsapi.OnMessageCb"))

	ws, err := NewWebSocketApi(mockProvider)
	require.NoError(t, err)

	ws.OnMessage(cb)
	mockProvider.AssertCalled(t, "OnMessage", mock.AnythingOfType("wsapi.OnMessageCb"))
}

func TestWebSocketApi_Send(t *testing.T) {
	mockProvider := new(MockWsProviderConnector)
	message := map[string]string{"test": "data"}
	mockProvider.Mock.On("Send", message).Return(nil)

	ws, err := NewWebSocketApi(mockProvider)
	require.NoError(t, err)

	err = ws.Send(message)
	assert.NoError(t, err)
	mockProvider.AssertCalled(t, "Send", message)
}

func TestWebSocketApi_Send_Error(t *testing.T) {
	mockProvider := new(MockWsProviderConnector)
	message := map[string]string{"test": "data"}
	mockProvider.Mock.On("Send", message).Return(errors.New("send error"))

	ws, err := NewWebSocketApi(mockProvider)
	require.NoError(t, err)

	err = ws.Send(message)
	assert.Error(t, err)
	assert.Equal(t, "send error", err.Error())
}

func TestWebSocketApi_Close(t *testing.T) {
	mockProvider := new(MockWsProviderConnector)
	mockProvider.Mock.On("Close").Return(nil)

	ws, err := NewWebSocketApi(mockProvider)
	require.NoError(t, err)

	err = ws.Close()
	assert.NoError(t, err)
	mockProvider.AssertCalled(t, "Close")
}

func TestWebSocketApi_Close_Error(t *testing.T) {
	mockProvider := new(MockWsProviderConnector)
	mockProvider.Mock.On("Close").Return(errors.New("close error"))

	ws, err := NewWebSocketApi(mockProvider)
	require.NoError(t, err)

	err = ws.Close()
	assert.Error(t, err)
	assert.Equal(t, "close error", err.Error())
}

func TestCastURL(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "wss URL",
			input:    "wss://api.1inch.dev/ws",
			expected: "wss://api.1inch.dev/ws",
		},
		{
			name:     "ws URL",
			input:    "ws://api.1inch.dev/ws",
			expected: "ws://api.1inch.dev/ws",
		},
		{
			name:     "https URL",
			input:    "https://api.1inch.dev/ws",
			expected: "wss://api.1inch.dev/ws",
		},
		{
			name:     "http URL",
			input:    "http://api.1inch.dev/ws",
			expected: "ws://api.1inch.dev/ws",
		},
		{
			name:     "no protocol",
			input:    "api.1inch.dev/ws",
			expected: "wss://api.1inch.dev/ws",
		},
		{
			name:     "URL with trailing slash",
			input:    "wss://api.1inch.dev/ws/",
			expected: "wss://api.1inch.dev/ws",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := castURL(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
