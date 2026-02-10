package wsapi

import (
	"testing"

	"github.com/dawitel/cross-chain-sdk-go/chains"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestNewActiveOrdersWebSocketApi(t *testing.T) {
	mockProvider := new(MockWsProviderConnector)
	api := NewActiveOrdersWebSocketApi(mockProvider)

	assert.NotNil(t, api)
	assert.Equal(t, mockProvider, api.provider)
}

func TestActiveOrdersWebSocketApi_OnOrder(t *testing.T) {
	mockProvider := new(MockWsProviderConnector)
	api := NewActiveOrdersWebSocketApi(mockProvider)

	var receivedEvent OrderEventType
	cb := func(data OrderEventType) {
		receivedEvent = data
	}

	// Set up mock to capture OnMessage callback
	var messageCb OnMessageCb
	mockProvider.Mock.On("OnMessage", mock.Anything).Run(func(args mock.Arguments) {
		if cb, ok := args.Get(0).(OnMessageCb); ok {
			messageCb = cb
		}
	})

	api.OnOrder(cb)
	mockProvider.AssertCalled(t, "OnMessage", mock.Anything)

	// Test with order_created event
	orderCreatedData := map[string]interface{}{
		"event":        "order_created",
		"orderHash":    "0x123",
		"srcChainId":   float64(chains.Ethereum),
		"dstChainId":   float64(chains.Polygon),
		"order":        map[string]interface{}{},
		"extension":    "",
		"signature":    "",
		"quoteId":      "quote-123",
		"merkleLeaves": []interface{}{},
		"secretHashes": []interface{}{},
	}

	if messageCb != nil {
		messageCb(orderCreatedData)
		assert.NotNil(t, receivedEvent)
		assert.Equal(t, EventTypeOrderCreated, receivedEvent.GetEvent())
	}
}

func TestActiveOrdersWebSocketApi_OnOrderCreated(t *testing.T) {
	mockProvider := new(MockWsProviderConnector)
	api := NewActiveOrdersWebSocketApi(mockProvider)

	var receivedEvent *OrderCreatedEvent
	cb := func(data *OrderCreatedEvent) {
		receivedEvent = data
	}

	var messageCb OnMessageCb
	mockProvider.Mock.On("OnMessage", mock.Anything).Run(func(args mock.Arguments) {
		if cb, ok := args.Get(0).(OnMessageCb); ok {
			messageCb = cb
		}
	})

	api.OnOrderCreated(cb)
	mockProvider.AssertCalled(t, "OnMessage", mock.Anything)

	orderCreatedData := map[string]interface{}{
		"event":        "order_created",
		"orderHash":    "0x123",
		"srcChainId":   float64(chains.Ethereum),
		"dstChainId":   float64(chains.Polygon),
		"order":        map[string]interface{}{},
		"extension":    "",
		"signature":    "",
		"quoteId":      "quote-123",
		"merkleLeaves": []interface{}{},
		"secretHashes": []interface{}{},
	}

	messageCb(orderCreatedData)
	assert.NotNil(t, receivedEvent)
	assert.Equal(t, "0x123", receivedEvent.OrderHash)
}

func TestActiveOrdersWebSocketApi_OnOrderInvalid(t *testing.T) {
	mockProvider := new(MockWsProviderConnector)
	api := NewActiveOrdersWebSocketApi(mockProvider)

	var receivedEvent *OrderInvalidEvent
	cb := func(data *OrderInvalidEvent) {
		receivedEvent = data
	}

	var messageCb OnMessageCb
	mockProvider.Mock.On("OnMessage", mock.Anything).Run(func(args mock.Arguments) {
		if cb, ok := args.Get(0).(OnMessageCb); ok {
			messageCb = cb
		}
	})

	api.OnOrderInvalid(cb)

	orderInvalidData := map[string]interface{}{
		"event":     "order_invalid",
		"orderHash": "0x456",
		"reason":    "insufficient balance",
	}

	messageCb(orderInvalidData)
	assert.NotNil(t, receivedEvent)
	assert.Equal(t, "0x456", receivedEvent.OrderHash)
	assert.Equal(t, "insufficient balance", receivedEvent.Reason)
}

func TestActiveOrdersWebSocketApi_OnOrderFilled(t *testing.T) {
	mockProvider := new(MockWsProviderConnector)
	api := NewActiveOrdersWebSocketApi(mockProvider)

	var receivedEvent *OrderFilledEvent
	cb := func(data *OrderFilledEvent) {
		receivedEvent = data
	}

	var messageCb OnMessageCb
	mockProvider.Mock.On("OnMessage", mock.Anything).Run(func(args mock.Arguments) {
		if cb, ok := args.Get(0).(OnMessageCb); ok {
			messageCb = cb
		}
	})

	api.OnOrderFilled(cb)

	orderFilledData := map[string]interface{}{
		"event":     "order_filled",
		"orderHash": "0x789",
		"txHash":    "0xtx123",
	}

	messageCb(orderFilledData)
	assert.NotNil(t, receivedEvent)
	assert.Equal(t, "0x789", receivedEvent.OrderHash)
	assert.Equal(t, "0xtx123", receivedEvent.TxHash)
}

func TestActiveOrdersWebSocketApi_OnOrderCancelled(t *testing.T) {
	mockProvider := new(MockWsProviderConnector)
	api := NewActiveOrdersWebSocketApi(mockProvider)

	var receivedEvent *OrderCancelledEvent
	cb := func(data *OrderCancelledEvent) {
		receivedEvent = data
	}

	var messageCb OnMessageCb
	mockProvider.Mock.On("OnMessage", mock.Anything).Run(func(args mock.Arguments) {
		if cb, ok := args.Get(0).(OnMessageCb); ok {
			messageCb = cb
		}
	})

	api.OnOrderCancelled(cb)

	orderCancelledData := map[string]interface{}{
		"event":                "order_cancelled",
		"orderHash":            "0xabc",
		"remainingMakerAmount": "1000000",
	}

	messageCb(orderCancelledData)
	assert.NotNil(t, receivedEvent)
	assert.Equal(t, "0xabc", receivedEvent.OrderHash)
	assert.Equal(t, "1000000", receivedEvent.RemainingMakerAmount)
}

func TestActiveOrdersWebSocketApi_OnOrderBalanceChange(t *testing.T) {
	mockProvider := new(MockWsProviderConnector)
	api := NewActiveOrdersWebSocketApi(mockProvider)

	var receivedEvent *OrderBalanceChangeEvent
	cb := func(data *OrderBalanceChangeEvent) {
		receivedEvent = data
	}

	var messageCb OnMessageCb
	mockProvider.Mock.On("OnMessage", mock.Anything).Run(func(args mock.Arguments) {
		if cb, ok := args.Get(0).(OnMessageCb); ok {
			messageCb = cb
		}
	})

	api.OnOrderBalanceChange(cb)

	balanceChangeData := map[string]interface{}{
		"event":                "order_balance_change",
		"orderHash":            "0xdef",
		"remainingMakerAmount": "500000",
		"balance":              "1000000",
	}

	messageCb(balanceChangeData)
	assert.NotNil(t, receivedEvent)
	assert.Equal(t, "0xdef", receivedEvent.OrderHash)
	assert.Equal(t, "500000", receivedEvent.RemainingMakerAmount)
	assert.Equal(t, "1000000", receivedEvent.Balance)
}

func TestActiveOrdersWebSocketApi_OnOrderAllowanceChange(t *testing.T) {
	mockProvider := new(MockWsProviderConnector)
	api := NewActiveOrdersWebSocketApi(mockProvider)

	var receivedEvent *OrderAllowanceChangeEvent
	cb := func(data *OrderAllowanceChangeEvent) {
		receivedEvent = data
	}

	var messageCb OnMessageCb
	mockProvider.Mock.On("OnMessage", mock.Anything).Run(func(args mock.Arguments) {
		if cb, ok := args.Get(0).(OnMessageCb); ok {
			messageCb = cb
		}
	})

	api.OnOrderAllowanceChange(cb)

	allowanceChangeData := map[string]interface{}{
		"event":                "order_allowance_change",
		"orderHash":            "0xghi",
		"remainingMakerAmount": "300000",
		"allowance":            "500000",
	}

	messageCb(allowanceChangeData)
	assert.NotNil(t, receivedEvent)
	assert.Equal(t, "0xghi", receivedEvent.OrderHash)
	assert.Equal(t, "300000", receivedEvent.RemainingMakerAmount)
	assert.Equal(t, "500000", receivedEvent.Allowance)
}

func TestActiveOrdersWebSocketApi_OnOrderFilledPartially(t *testing.T) {
	mockProvider := new(MockWsProviderConnector)
	api := NewActiveOrdersWebSocketApi(mockProvider)

	var receivedEvent *OrderFilledPartiallyEvent
	cb := func(data *OrderFilledPartiallyEvent) {
		receivedEvent = data
	}

	var messageCb OnMessageCb
	mockProvider.Mock.On("OnMessage", mock.Anything).Run(func(args mock.Arguments) {
		if cb, ok := args.Get(0).(OnMessageCb); ok {
			messageCb = cb
		}
	})

	api.OnOrderFilledPartially(cb)

	partialFilledData := map[string]interface{}{
		"event":     "order_filled_partially",
		"orderHash": "0xjkl",
		"txHash":    "0xtx456",
	}

	messageCb(partialFilledData)
	assert.NotNil(t, receivedEvent)
	assert.Equal(t, "0xjkl", receivedEvent.OrderHash)
	assert.Equal(t, "0xtx456", receivedEvent.TxHash)
}

func TestActiveOrdersWebSocketApi_OnOrderSecretShared(t *testing.T) {
	mockProvider := new(MockWsProviderConnector)
	api := NewActiveOrdersWebSocketApi(mockProvider)

	var receivedEvent *OrderSecretSharedEvent
	cb := func(data *OrderSecretSharedEvent) {
		receivedEvent = data
	}

	var messageCb OnMessageCb
	mockProvider.Mock.On("OnMessage", mock.Anything).Run(func(args mock.Arguments) {
		if cb, ok := args.Get(0).(OnMessageCb); ok {
			messageCb = cb
		}
	})

	api.OnOrderSecretShared(cb)

	secretSharedData := map[string]interface{}{
		"event":  "secret_shared",
		"idx":    float64(0),
		"secret": "0xsecret123",
	}

	messageCb(secretSharedData)
	assert.NotNil(t, receivedEvent)
	assert.Equal(t, 0, receivedEvent.Idx)
	assert.Equal(t, "0xsecret123", receivedEvent.Secret)
}

func TestActiveOrdersWebSocketApi_parseOrderEvent_InvalidData(t *testing.T) {
	mockProvider := new(MockWsProviderConnector)
	api := NewActiveOrdersWebSocketApi(mockProvider)

	// Test with invalid JSON-like data
	invalidData := "not a map"
	event, ok := api.parseOrderEvent(invalidData)
	assert.False(t, ok)
	assert.Nil(t, event)

	// Test with missing event field
	invalidData2 := map[string]interface{}{
		"orderHash": "0x123",
	}
	event, ok = api.parseOrderEvent(invalidData2)
	assert.False(t, ok)
	assert.Nil(t, event)

	// Test with unknown event type
	unknownEventData := map[string]interface{}{
		"event": "unknown_event",
	}
	event, ok = api.parseOrderEvent(unknownEventData)
	assert.False(t, ok)
	assert.Nil(t, event)
}

func TestActiveOrdersWebSocketApi_parseOrderEvent_OrderCreated_WithOrder(t *testing.T) {
	mockProvider := new(MockWsProviderConnector)
	api := NewActiveOrdersWebSocketApi(mockProvider)

	orderData := map[string]interface{}{
		"maker":        "0xmaker",
		"makerAsset":   "0xasset",
		"takerAsset":   "0xtaker",
		"makingAmount": "1000000",
		"takingAmount": "950000",
		"receiver":     "0xreceiver",
		"salt":         "12345",
		"expiration":   "1735689600",
		"nonce":        "67890",
	}

	orderCreatedData := map[string]interface{}{
		"event":        "order_created",
		"orderHash":    "0x123",
		"srcChainId":   float64(chains.Ethereum),
		"dstChainId":   float64(chains.Polygon),
		"order":        orderData,
		"extension":    "0xext",
		"signature":    "0xsig",
		"quoteId":      "quote-123",
		"merkleLeaves": []interface{}{"0xleaf1"},
		"secretHashes": []interface{}{"0xhash1"},
	}

	event, ok := api.parseOrderEvent(orderCreatedData)
	require.True(t, ok)
	assert.NotNil(t, event)
	assert.Equal(t, EventTypeOrderCreated, event.GetEvent())

	createdEvent, ok := event.(*OrderCreatedEvent)
	require.True(t, ok)
	assert.Equal(t, "0x123", createdEvent.OrderHash)
	assert.Equal(t, chains.Ethereum, createdEvent.SrcChainID)
	assert.Equal(t, chains.Polygon, createdEvent.DstChainID)
	assert.Equal(t, "0xext", createdEvent.Extension)
	assert.Equal(t, "0xsig", createdEvent.Signature)
	assert.Equal(t, "quote-123", createdEvent.QuoteID)
	assert.Len(t, createdEvent.MerkleLeaves, 1)
	assert.Len(t, createdEvent.SecretHashes, 1)
}
