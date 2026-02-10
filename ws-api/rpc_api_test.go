package wsapi

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNewRpcWebsocketApi(t *testing.T) {
	mockProvider := new(MockWsProviderConnector)
	api := NewRpcWebsocketApi(mockProvider)

	assert.NotNil(t, api)
	assert.Equal(t, mockProvider, api.provider)
}

func TestRpcWebsocketApi_Ping(t *testing.T) {
	mockProvider := new(MockWsProviderConnector)
	api := NewRpcWebsocketApi(mockProvider)

	expectedMessage := map[string]interface{}{
		"method": RpcMethodPing,
	}
	mockProvider.Mock.On("Send", expectedMessage).Return(nil)

	err := api.Ping()
	assert.NoError(t, err)
	mockProvider.AssertCalled(t, "Send", expectedMessage)
}

func TestRpcWebsocketApi_Ping_Error(t *testing.T) {
	mockProvider := new(MockWsProviderConnector)
	api := NewRpcWebsocketApi(mockProvider)

	mockProvider.Mock.On("Send", mock.Anything).Return(errors.New("send error"))

	err := api.Ping()
	assert.Error(t, err)
	assert.Equal(t, "send error", err.Error())
}

func TestRpcWebsocketApi_OnPong(t *testing.T) {
	mockProvider := new(MockWsProviderConnector)
	api := NewRpcWebsocketApi(mockProvider)

	var receivedResult string
	cb := func(result string) {
		receivedResult = result
	}

	var messageCb OnMessageCb
	mockProvider.Mock.On("OnMessage", mock.Anything).Run(func(args mock.Arguments) {
		if cb, ok := args.Get(0).(OnMessageCb); ok {
			messageCb = cb
		}
	})

	api.OnPong(cb)
	mockProvider.AssertCalled(t, "OnMessage", mock.Anything)

	pingData := map[string]interface{}{
		"method": "ping",
		"result": "pong",
	}

	messageCb(pingData)
	assert.Equal(t, "pong", receivedResult)
}

func TestRpcWebsocketApi_GetActiveOrders(t *testing.T) {
	mockProvider := new(MockWsProviderConnector)
	api := NewRpcWebsocketApi(mockProvider)

	expectedMessage := map[string]interface{}{
		"method": RpcMethodGetActiveOrders,
		"param": map[string]interface{}{
			"page":  1,
			"limit": 10,
		},
	}
	mockProvider.Mock.On("Send", expectedMessage).Return(nil)

	err := api.GetActiveOrders(1, 10)
	assert.NoError(t, err)
	mockProvider.AssertCalled(t, "Send", expectedMessage)
}

func TestRpcWebsocketApi_GetActiveOrders_ZeroParams(t *testing.T) {
	mockProvider := new(MockWsProviderConnector)
	api := NewRpcWebsocketApi(mockProvider)

	expectedMessage := map[string]interface{}{
		"method": RpcMethodGetActiveOrders,
		"param":  map[string]interface{}{},
	}
	mockProvider.Mock.On("Send", expectedMessage).Return(nil)

	err := api.GetActiveOrders(0, 0)
	assert.NoError(t, err)
	mockProvider.AssertCalled(t, "Send", expectedMessage)
}

func TestRpcWebsocketApi_OnGetActiveOrders(t *testing.T) {
	mockProvider := new(MockWsProviderConnector)
	api := NewRpcWebsocketApi(mockProvider)

	var receivedResult interface{}
	cb := func(result interface{}) {
		receivedResult = result
	}

	var messageCb OnMessageCb
	mockProvider.Mock.On("OnMessage", mock.Anything).Run(func(args mock.Arguments) {
		if cb, ok := args.Get(0).(OnMessageCb); ok {
			messageCb = cb
		}
	})

	api.OnGetActiveOrders(cb)

	resultData := map[string]interface{}{
		"totalCount": float64(5),
		"items":      []interface{}{},
	}

	ordersData := map[string]interface{}{
		"method": "getActiveOrders",
		"result": resultData,
	}

	messageCb(ordersData)
	assert.NotNil(t, receivedResult)
}

func TestRpcWebsocketApi_GetSecrets(t *testing.T) {
	mockProvider := new(MockWsProviderConnector)
	api := NewRpcWebsocketApi(mockProvider)

	expectedMessage := map[string]interface{}{
		"method": RpcMethodGetSecrets,
		"param": map[string]interface{}{
			"page":  2,
			"limit": 20,
		},
	}
	mockProvider.Mock.On("Send", expectedMessage).Return(nil)

	err := api.GetSecrets(2, 20)
	assert.NoError(t, err)
	mockProvider.AssertCalled(t, "Send", expectedMessage)
}

func TestRpcWebsocketApi_OnGetSecrets(t *testing.T) {
	mockProvider := new(MockWsProviderConnector)
	api := NewRpcWebsocketApi(mockProvider)

	var receivedResult interface{}
	cb := func(result interface{}) {
		receivedResult = result
	}

	var messageCb OnMessageCb
	mockProvider.Mock.On("OnMessage", mock.Anything).Run(func(args mock.Arguments) {
		if cb, ok := args.Get(0).(OnMessageCb); ok {
			messageCb = cb
		}
	})

	api.OnGetSecrets(cb)

	resultData := map[string]interface{}{
		"secrets": []interface{}{"0xsecret1", "0xsecret2"},
	}

	secretsData := map[string]interface{}{
		"method": "getSecrets",
		"result": resultData,
	}

	messageCb(secretsData)
	assert.NotNil(t, receivedResult)
}

func TestRpcWebsocketApi_GetAllowedMethods(t *testing.T) {
	mockProvider := new(MockWsProviderConnector)
	api := NewRpcWebsocketApi(mockProvider)

	expectedMessage := map[string]interface{}{
		"method": RpcMethodGetAllowedMethods,
	}
	mockProvider.Mock.On("Send", expectedMessage).Return(nil)

	err := api.GetAllowedMethods()
	assert.NoError(t, err)
	mockProvider.AssertCalled(t, "Send", expectedMessage)
}

func TestRpcWebsocketApi_OnGetAllowedMethods(t *testing.T) {
	mockProvider := new(MockWsProviderConnector)
	api := NewRpcWebsocketApi(mockProvider)

	var receivedResult []string
	cb := func(result []string) {
		receivedResult = result
	}

	var messageCb OnMessageCb
	mockProvider.Mock.On("OnMessage", mock.Anything).Run(func(args mock.Arguments) {
		if cb, ok := args.Get(0).(OnMessageCb); ok {
			messageCb = cb
		}
	})

	api.OnGetAllowedMethods(cb)

	methodsData := map[string]interface{}{
		"method": "getAllowedMethods",
		"result": []interface{}{"ping", "getActiveOrders", "getSecrets"},
	}

	messageCb(methodsData)
	assert.NotNil(t, receivedResult)
	assert.Len(t, receivedResult, 3)
	assert.Contains(t, receivedResult, "ping")
	assert.Contains(t, receivedResult, "getActiveOrders")
	assert.Contains(t, receivedResult, "getSecrets")
}

func TestRpcWebsocketApi_parseRpcEvent_InvalidData(t *testing.T) {
	mockProvider := new(MockWsProviderConnector)
	api := NewRpcWebsocketApi(mockProvider)

	// Test with invalid data
	invalidData := "not a map"
	event, ok := api.parseRpcEvent(invalidData)
	assert.False(t, ok)
	assert.Nil(t, event)

	// Test with missing method field
	invalidData2 := map[string]interface{}{
		"result": "some result",
	}
	event, ok = api.parseRpcEvent(invalidData2)
	assert.False(t, ok)
	assert.Nil(t, event)

	// Test with unknown method
	unknownMethodData := map[string]interface{}{
		"method": "unknown_method",
		"result": "some result",
	}
	event, ok = api.parseRpcEvent(unknownMethodData)
	assert.False(t, ok)
	assert.Nil(t, event)
}

func TestRpcWebsocketApi_parseRpcEvent_AllMethods(t *testing.T) {
	mockProvider := new(MockWsProviderConnector)
	api := NewRpcWebsocketApi(mockProvider)

	tests := []struct {
		name     string
		method   RpcMethod
		data     map[string]interface{}
		validate func(t *testing.T, event RpcEventType, ok bool)
	}{
		{
			name:   "ping",
			method: RpcMethodPing,
			data: map[string]interface{}{
				"method": "ping",
				"result": "pong",
			},
			validate: func(t *testing.T, event RpcEventType, ok bool) {
				require.True(t, ok)
				assert.Equal(t, RpcMethodPing, event.GetMethod())
				pingEvent, ok := event.(*PingRpcEvent)
				require.True(t, ok)
				assert.Equal(t, "pong", pingEvent.Result)
			},
		},
		{
			name:   "getActiveOrders",
			method: RpcMethodGetActiveOrders,
			data: map[string]interface{}{
				"method": "getActiveOrders",
				"result": map[string]interface{}{"totalCount": float64(10)},
			},
			validate: func(t *testing.T, event RpcEventType, ok bool) {
				require.True(t, ok)
				assert.Equal(t, RpcMethodGetActiveOrders, event.GetMethod())
			},
		},
		{
			name:   "getSecrets",
			method: RpcMethodGetSecrets,
			data: map[string]interface{}{
				"method": "getSecrets",
				"result": map[string]interface{}{"secrets": []interface{}{}},
			},
			validate: func(t *testing.T, event RpcEventType, ok bool) {
				require.True(t, ok)
				assert.Equal(t, RpcMethodGetSecrets, event.GetMethod())
			},
		},
		{
			name:   "getAllowedMethods",
			method: RpcMethodGetAllowedMethods,
			data: map[string]interface{}{
				"method": "getAllowedMethods",
				"result": []interface{}{"ping"},
			},
			validate: func(t *testing.T, event RpcEventType, ok bool) {
				require.True(t, ok)
				assert.Equal(t, RpcMethodGetAllowedMethods, event.GetMethod())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event, ok := api.parseRpcEvent(tt.data)
			tt.validate(t, event, ok)
		})
	}
}

func TestRpcWebsocketApi_OnPong_WithInvalidData(t *testing.T) {
	mockProvider := new(MockWsProviderConnector)
	api := NewRpcWebsocketApi(mockProvider)

	cbCalled := false
	cb := func(result string) {
		cbCalled = true
	}

	var messageCb OnMessageCb
	mockProvider.Mock.On("OnMessage", mock.Anything).Run(func(args mock.Arguments) {
		if cb, ok := args.Get(0).(OnMessageCb); ok {
			messageCb = cb
		}
	})

	api.OnPong(cb)

	// Send non-ping data
	if messageCb != nil {
		messageCb(map[string]interface{}{"method": "getActiveOrders"})
		assert.False(t, cbCalled)

		// Send ping data without result
		messageCb(map[string]interface{}{"method": "ping"})
		assert.False(t, cbCalled)
	}
}
