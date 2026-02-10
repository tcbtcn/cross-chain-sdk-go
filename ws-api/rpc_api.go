package wsapi

import (
	"encoding/json"
)

// RpcWebsocketApi handles RPC-related WebSocket methods
type RpcWebsocketApi struct {
	provider WsProviderConnector
}

// NewRpcWebsocketApi creates a new RpcWebsocketApi
func NewRpcWebsocketApi(provider WsProviderConnector) *RpcWebsocketApi {
	return &RpcWebsocketApi{
		provider: provider,
	}
}

// Ping sends a ping message
func (r *RpcWebsocketApi) Ping() error {
	return r.provider.Send(map[string]interface{}{
		"method": RpcMethodPing,
	})
}

// OnPong subscribes to pong responses
func (r *RpcWebsocketApi) OnPong(cb OnPongCb) {
	r.provider.OnMessage(func(data interface{}) {
		rpcEvent, ok := r.parseRpcEvent(data)
		if ok && rpcEvent.GetMethod() == RpcMethodPing {
			if pingEvent, ok := rpcEvent.(*PingRpcEvent); ok {
				if pingEvent.Result != "" {
					cb(pingEvent.Result)
				}
			}
		}
	})
}

// GetActiveOrders requests active orders
func (r *RpcWebsocketApi) GetActiveOrders(page, limit int) error {
	params := map[string]interface{}{}
	if page > 0 {
		params["page"] = page
	}
	if limit > 0 {
		params["limit"] = limit
	}

	return r.provider.Send(map[string]interface{}{
		"method": RpcMethodGetActiveOrders,
		"param":  params,
	})
}

// OnGetActiveOrders subscribes to getActiveOrders responses
func (r *RpcWebsocketApi) OnGetActiveOrders(cb OnGetActiveOrdersCb) {
	r.provider.OnMessage(func(data interface{}) {
		rpcEvent, ok := r.parseRpcEvent(data)
		if ok && rpcEvent.GetMethod() == RpcMethodGetActiveOrders {
			if ordersEvent, ok := rpcEvent.(*GetActiveOrdersRpcEvent); ok {
				cb(ordersEvent.Result)
			}
		}
	})
}

// GetSecrets requests secrets
func (r *RpcWebsocketApi) GetSecrets(page, limit int) error {
	params := map[string]interface{}{}
	if page > 0 {
		params["page"] = page
	}
	if limit > 0 {
		params["limit"] = limit
	}

	return r.provider.Send(map[string]interface{}{
		"method": RpcMethodGetSecrets,
		"param":  params,
	})
}

// OnGetSecrets subscribes to getSecrets responses
func (r *RpcWebsocketApi) OnGetSecrets(cb OnGetSecretsCb) {
	r.provider.OnMessage(func(data interface{}) {
		rpcEvent, ok := r.parseRpcEvent(data)
		if ok && rpcEvent.GetMethod() == RpcMethodGetSecrets {
			if secretsEvent, ok := rpcEvent.(*GetSecretsRpcEvent); ok {
				cb(secretsEvent.Result)
			}
		}
	})
}

// GetAllowedMethods requests allowed methods
func (r *RpcWebsocketApi) GetAllowedMethods() error {
	return r.provider.Send(map[string]interface{}{
		"method": RpcMethodGetAllowedMethods,
	})
}

// OnGetAllowedMethods subscribes to getAllowedMethods responses
func (r *RpcWebsocketApi) OnGetAllowedMethods(cb OnGetAllowedMethodsCb) {
	r.provider.OnMessage(func(data interface{}) {
		rpcEvent, ok := r.parseRpcEvent(data)
		if ok && rpcEvent.GetMethod() == RpcMethodGetAllowedMethods {
			if methodsEvent, ok := rpcEvent.(*GetAllowedMethodsRpcEvent); ok {
				cb(methodsEvent.Result)
			}
		}
	})
}

func (r *RpcWebsocketApi) parseRpcEvent(data interface{}) (RpcEventType, bool) {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return nil, false
	}

	var eventMap map[string]interface{}
	if err := json.Unmarshal(dataBytes, &eventMap); err != nil {
		return nil, false
	}

	methodStr, ok := eventMap["method"].(string)
	if !ok {
		return nil, false
	}

	method := RpcMethod(methodStr)

	switch method {
	case RpcMethodPing:
		var event PingRpcEvent
		if err := json.Unmarshal(dataBytes, &event); err == nil {
			return &event, true
		}
	case RpcMethodGetActiveOrders:
		var event GetActiveOrdersRpcEvent
		if err := json.Unmarshal(dataBytes, &event); err == nil {
			return &event, true
		}
	case RpcMethodGetSecrets:
		var event GetSecretsRpcEvent
		if err := json.Unmarshal(dataBytes, &event); err == nil {
			return &event, true
		}
	case RpcMethodGetAllowedMethods:
		var event GetAllowedMethodsRpcEvent
		if err := json.Unmarshal(dataBytes, &event); err == nil {
			return &event, true
		}
	}

	return nil, false
}
