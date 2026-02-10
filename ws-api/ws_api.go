package wsapi

import (
	"fmt"
	"strings"
)

const Version = "v1.0"

// WebSocketApi is the main WebSocket API client
type WebSocketApi struct {
	RPC      *RpcWebsocketApi
	Order    *ActiveOrdersWebSocketApi
	provider WsProviderConnector
}

// NewWebSocketApi creates a new WebSocketApi
func NewWebSocketApi(configOrProvider interface{}) (*WebSocketApi, error) {
	var provider WsProviderConnector

	switch v := configOrProvider.(type) {
	case WsProviderConnector:
		provider = v
	case WsApiConfig:
		url := castURL(v.URL)
		config := WsApiConfig{
			URL:      fmt.Sprintf("%s/%s", url, Version),
			AuthKey:  v.AuthKey,
			LazyInit: v.LazyInit,
		}
		provider = NewWebSocketClient(config)
	default:
		return nil, fmt.Errorf("invalid config or provider type")
	}

	return &WebSocketApi{
		RPC:      NewRpcWebsocketApi(provider),
		Order:    NewActiveOrdersWebSocketApi(provider),
		provider: provider,
	}, nil
}

// Init initializes the WebSocket connection
func (w *WebSocketApi) Init() {
	w.provider.Init()
}

// On subscribes to a WebSocket event
func (w *WebSocketApi) On(event WebSocketEvent, cb interface{}) {
	w.provider.On(event, cb)
}

// Off unsubscribes from a WebSocket event
func (w *WebSocketApi) Off(event WebSocketEvent, cb interface{}) {
	w.provider.Off(event, cb)
}

// OnOpen subscribes to the open event
func (w *WebSocketApi) OnOpen(cb OnOpenCb) {
	w.provider.OnOpen(cb)
}

// OnClose subscribes to the close event
func (w *WebSocketApi) OnClose(cb OnCloseCb) {
	w.provider.OnClose(cb)
}

// OnError subscribes to the error event
func (w *WebSocketApi) OnError(cb OnErrorCb) {
	w.provider.OnError(cb)
}

// OnMessage subscribes to the message event
func (w *WebSocketApi) OnMessage(cb OnMessageCb) {
	w.provider.OnMessage(cb)
}

// Send sends a message through the WebSocket
func (w *WebSocketApi) Send(message interface{}) error {
	return w.provider.Send(message)
}

// Close closes the WebSocket connection
func (w *WebSocketApi) Close() error {
	return w.provider.Close()
}

func castURL(url string) string {
	url = strings.TrimSuffix(url, "/")
	if !strings.HasPrefix(url, "ws://") && !strings.HasPrefix(url, "wss://") {
		if strings.HasPrefix(url, "http://") {
			url = strings.Replace(url, "http://", "ws://", 1)
		} else if strings.HasPrefix(url, "https://") {
			url = strings.Replace(url, "https://", "wss://", 1)
		} else {
			url = "wss://" + url
		}
	}
	return url
}
