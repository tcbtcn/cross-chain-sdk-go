package wsapi

import (
	"testing"

	"github.com/dawitel/cross-chain-sdk-go/chains"
	"github.com/dawitel/cross-chain-sdk-go/orders"
	"github.com/stretchr/testify/assert"
)

func TestWebSocketEvent(t *testing.T) {
	assert.Equal(t, WebSocketEvent("close"), EventClose)
	assert.Equal(t, WebSocketEvent("error"), EventError)
	assert.Equal(t, WebSocketEvent("message"), EventMessage)
	assert.Equal(t, WebSocketEvent("open"), EventOpen)
}

func TestEventType(t *testing.T) {
	assert.Equal(t, EventType("order_created"), EventTypeOrderCreated)
	assert.Equal(t, EventType("order_invalid"), EventTypeOrderInvalid)
	assert.Equal(t, EventType("order_balance_change"), EventTypeOrderBalanceChange)
	assert.Equal(t, EventType("order_allowance_change"), EventTypeOrderAllowanceChange)
	assert.Equal(t, EventType("order_filled"), EventTypeOrderFilled)
	assert.Equal(t, EventType("order_filled_partially"), EventTypeOrderFilledPartially)
	assert.Equal(t, EventType("order_cancelled"), EventTypeOrderCancelled)
	assert.Equal(t, EventType("secret_shared"), EventTypeOrderSecretShared)
}

func TestRpcMethod(t *testing.T) {
	assert.Equal(t, RpcMethod("getAllowedMethods"), RpcMethodGetAllowedMethods)
	assert.Equal(t, RpcMethod("ping"), RpcMethodPing)
	assert.Equal(t, RpcMethod("getActiveOrders"), RpcMethodGetActiveOrders)
	assert.Equal(t, RpcMethod("getSecrets"), RpcMethodGetSecrets)
}

func TestOrderCreatedEvent_GetEvent(t *testing.T) {
	event := &OrderCreatedEvent{
		Event:     EventTypeOrderCreated,
		OrderHash: "0x123",
	}
	assert.Equal(t, EventTypeOrderCreated, event.GetEvent())
}

func TestOrderInvalidEvent_GetEvent(t *testing.T) {
	event := &OrderInvalidEvent{
		Event:     EventTypeOrderInvalid,
		OrderHash: "0x456",
	}
	assert.Equal(t, EventTypeOrderInvalid, event.GetEvent())
}

func TestOrderBalanceChangeEvent_GetEvent(t *testing.T) {
	event := &OrderBalanceChangeEvent{
		Event:     EventTypeOrderBalanceChange,
		OrderHash: "0x789",
	}
	assert.Equal(t, EventTypeOrderBalanceChange, event.GetEvent())
}

func TestOrderAllowanceChangeEvent_GetEvent(t *testing.T) {
	event := &OrderAllowanceChangeEvent{
		Event:     EventTypeOrderAllowanceChange,
		OrderHash: "0xabc",
	}
	assert.Equal(t, EventTypeOrderAllowanceChange, event.GetEvent())
}

func TestOrderFilledEvent_GetEvent(t *testing.T) {
	event := &OrderFilledEvent{
		Event:     EventTypeOrderFilled,
		OrderHash: "0xdef",
	}
	assert.Equal(t, EventTypeOrderFilled, event.GetEvent())
}

func TestOrderFilledPartiallyEvent_GetEvent(t *testing.T) {
	event := &OrderFilledPartiallyEvent{
		Event:     EventTypeOrderFilledPartially,
		OrderHash: "0xghi",
	}
	assert.Equal(t, EventTypeOrderFilledPartially, event.GetEvent())
}

func TestOrderCancelledEvent_GetEvent(t *testing.T) {
	event := &OrderCancelledEvent{
		Event:     EventTypeOrderCancelled,
		OrderHash: "0xjkl",
	}
	assert.Equal(t, EventTypeOrderCancelled, event.GetEvent())
}

func TestOrderSecretSharedEvent_GetEvent(t *testing.T) {
	event := &OrderSecretSharedEvent{
		Event:  EventTypeOrderSecretShared,
		Idx:    0,
		Secret: "0xsecret",
	}
	assert.Equal(t, EventTypeOrderSecretShared, event.GetEvent())
}

func TestPingRpcEvent_GetMethod(t *testing.T) {
	event := &PingRpcEvent{
		Method: RpcMethodPing,
		Result: "pong",
	}
	assert.Equal(t, RpcMethodPing, event.GetMethod())
}

func TestGetActiveOrdersRpcEvent_GetMethod(t *testing.T) {
	event := &GetActiveOrdersRpcEvent{
		Method: RpcMethodGetActiveOrders,
		Result: map[string]interface{}{},
	}
	assert.Equal(t, RpcMethodGetActiveOrders, event.GetMethod())
}

func TestGetSecretsRpcEvent_GetMethod(t *testing.T) {
	event := &GetSecretsRpcEvent{
		Method: RpcMethodGetSecrets,
		Result: map[string]interface{}{},
	}
	assert.Equal(t, RpcMethodGetSecrets, event.GetMethod())
}

func TestGetAllowedMethodsRpcEvent_GetMethod(t *testing.T) {
	event := &GetAllowedMethodsRpcEvent{
		Method: RpcMethodGetAllowedMethods,
		Result: []string{"ping"},
	}
	assert.Equal(t, RpcMethodGetAllowedMethods, event.GetMethod())
}

func TestOrderCreatedEvent_Fields(t *testing.T) {
	order := orders.LimitOrderV4Struct{
		Maker:        "0xmaker",
		MakerAsset:   "0xasset",
		TakerAsset:   "0xtaker",
		MakingAmount: "1000000",
		TakingAmount: "950000",
		Receiver:     "0xreceiver",
		Salt:         "12345",
		Expiration:   "1735689600",
		Nonce:        "67890",
	}

	event := &OrderCreatedEvent{
		Event:           EventTypeOrderCreated,
		SrcChainID:      chains.Ethereum,
		DstChainID:      chains.Polygon,
		OrderHash:       "0xhash",
		Order:           order,
		Extension:       "0xext",
		Signature:       "0xsig",
		IsMakerContract: false,
		QuoteID:         "quote-123",
		MerkleLeaves:    []string{"0xleaf1", "0xleaf2"},
		SecretHashes:    []string{"0xhash1", "0xhash2"},
	}

	assert.Equal(t, EventTypeOrderCreated, event.Event)
	assert.Equal(t, chains.Ethereum, event.SrcChainID)
	assert.Equal(t, chains.Polygon, event.DstChainID)
	assert.Equal(t, "0xhash", event.OrderHash)
	assert.Equal(t, order, event.Order)
	assert.Equal(t, "0xext", event.Extension)
	assert.Equal(t, "0xsig", event.Signature)
	assert.False(t, event.IsMakerContract)
	assert.Equal(t, "quote-123", event.QuoteID)
	assert.Len(t, event.MerkleLeaves, 2)
	assert.Len(t, event.SecretHashes, 2)
}

func TestOrderEvents_AllTypes(t *testing.T) {
	events := []OrderEventType{
		&OrderCreatedEvent{Event: EventTypeOrderCreated},
		&OrderInvalidEvent{Event: EventTypeOrderInvalid},
		&OrderBalanceChangeEvent{Event: EventTypeOrderBalanceChange},
		&OrderAllowanceChangeEvent{Event: EventTypeOrderAllowanceChange},
		&OrderFilledEvent{Event: EventTypeOrderFilled},
		&OrderFilledPartiallyEvent{Event: EventTypeOrderFilledPartially},
		&OrderCancelledEvent{Event: EventTypeOrderCancelled},
		&OrderSecretSharedEvent{Event: EventTypeOrderSecretShared},
	}

	for _, event := range events {
		assert.NotNil(t, event.GetEvent())
	}
}

func TestRpcEvents_AllTypes(t *testing.T) {
	events := []RpcEventType{
		&PingRpcEvent{Method: RpcMethodPing},
		&GetActiveOrdersRpcEvent{Method: RpcMethodGetActiveOrders},
		&GetSecretsRpcEvent{Method: RpcMethodGetSecrets},
		&GetAllowedMethodsRpcEvent{Method: RpcMethodGetAllowedMethods},
	}

	for _, event := range events {
		assert.NotNil(t, event.GetMethod())
	}
}

func TestOrderEvents(t *testing.T) {
	// Verify OrderEvents constant contains all event types
	expectedEvents := []EventType{
		EventTypeOrderCreated,
		EventTypeOrderInvalid,
		EventTypeOrderBalanceChange,
		EventTypeOrderAllowanceChange,
		EventTypeOrderFilled,
		EventTypeOrderFilledPartially,
		EventTypeOrderCancelled,
		EventTypeOrderSecretShared,
	}

	assert.Len(t, OrderEvents, len(expectedEvents))
	for _, expected := range expectedEvents {
		found := false
		for _, actual := range OrderEvents {
			if actual == expected {
				found = true
				break
			}
		}
		assert.True(t, found, "Event type %s should be in OrderEvents", expected)
	}
}
