package wsapi

import (
	"github.com/dawitel/cross-chain-sdk-go/chains"
	"github.com/dawitel/cross-chain-sdk-go/orders"
)

// WebSocketEvent represents WebSocket connection events
type WebSocketEvent string

const (
	EventClose   WebSocketEvent = "close"
	EventError   WebSocketEvent = "error"
	EventMessage WebSocketEvent = "message"
	EventOpen    WebSocketEvent = "open"
)

// EventType represents order event types
type EventType string

const (
	EventTypeOrderCreated         EventType = "order_created"
	EventTypeOrderInvalid         EventType = "order_invalid"
	EventTypeOrderBalanceChange   EventType = "order_balance_change"
	EventTypeOrderAllowanceChange EventType = "order_allowance_change"
	EventTypeOrderFilled          EventType = "order_filled"
	EventTypeOrderFilledPartially EventType = "order_filled_partially"
	EventTypeOrderCancelled       EventType = "order_cancelled"
	EventTypeOrderSecretShared    EventType = "secret_shared"
)

// RpcMethod represents RPC method names
type RpcMethod string

const (
	RpcMethodGetAllowedMethods RpcMethod = "getAllowedMethods"
	RpcMethodPing              RpcMethod = "ping"
	RpcMethodGetActiveOrders   RpcMethod = "getActiveOrders"
	RpcMethodGetSecrets        RpcMethod = "getSecrets"
)

// OrderEventType represents any order event
type OrderEventType interface {
	GetEvent() EventType
}

// OrderCreatedEvent represents an order_created event
type OrderCreatedEvent struct {
	Event           EventType                 `json:"event"`
	SrcChainID      chains.SupportedChain     `json:"srcChainId"`
	DstChainID      chains.SupportedChain     `json:"dstChainId"`
	OrderHash       string                    `json:"orderHash"`
	Order           orders.LimitOrderV4Struct `json:"order"`
	Extension       string                    `json:"extension"`
	Signature       string                    `json:"signature"`
	IsMakerContract bool                      `json:"isMakerContract"`
	QuoteID         string                    `json:"quoteId"`
	MerkleLeaves    []string                  `json:"merkleLeaves"`
	SecretHashes    []string                  `json:"secretHashes"`
}

func (e *OrderCreatedEvent) GetEvent() EventType {
	return EventTypeOrderCreated
}

// OrderInvalidEvent represents an order_invalid event
type OrderInvalidEvent struct {
	Event     EventType `json:"event"`
	OrderHash string    `json:"orderHash"`
	Reason    string    `json:"reason,omitempty"`
}

func (e *OrderInvalidEvent) GetEvent() EventType {
	return EventTypeOrderInvalid
}

// OrderBalanceChangeEvent represents an order_balance_change event
type OrderBalanceChangeEvent struct {
	Event                EventType `json:"event"`
	OrderHash            string    `json:"orderHash"`
	RemainingMakerAmount string    `json:"remainingMakerAmount"`
	Balance              string    `json:"balance"`
}

func (e *OrderBalanceChangeEvent) GetEvent() EventType {
	return EventTypeOrderBalanceChange
}

// OrderAllowanceChangeEvent represents an order_allowance_change event
type OrderAllowanceChangeEvent struct {
	Event                EventType `json:"event"`
	OrderHash            string    `json:"orderHash"`
	RemainingMakerAmount string    `json:"remainingMakerAmount"`
	Allowance            string    `json:"allowance"`
}

func (e *OrderAllowanceChangeEvent) GetEvent() EventType {
	return EventTypeOrderAllowanceChange
}

// OrderFilledEvent represents an order_filled event
type OrderFilledEvent struct {
	Event     EventType `json:"event"`
	OrderHash string    `json:"orderHash"`
	TxHash    string    `json:"txHash,omitempty"`
}

func (e *OrderFilledEvent) GetEvent() EventType {
	return EventTypeOrderFilled
}

// OrderFilledPartiallyEvent represents an order_filled_partially event
type OrderFilledPartiallyEvent struct {
	Event     EventType `json:"event"`
	OrderHash string    `json:"orderHash"`
	TxHash    string    `json:"txHash,omitempty"`
}

func (e *OrderFilledPartiallyEvent) GetEvent() EventType {
	return EventTypeOrderFilledPartially
}

// OrderCancelledEvent represents an order_cancelled event
type OrderCancelledEvent struct {
	Event                EventType `json:"event"`
	OrderHash            string    `json:"orderHash"`
	RemainingMakerAmount string    `json:"remainingMakerAmount"`
}

func (e *OrderCancelledEvent) GetEvent() EventType {
	return EventTypeOrderCancelled
}

// OrderSecretSharedEvent represents a secret_shared event
type OrderSecretSharedEvent struct {
	Event         EventType   `json:"event"`
	Idx           int         `json:"idx"`
	Secret        string      `json:"secret"`
	SrcImmutables interface{} `json:"srcImmutables,omitempty"`
	DstImmutables interface{} `json:"dstImmutables,omitempty"`
}

func (e *OrderSecretSharedEvent) GetEvent() EventType {
	return EventTypeOrderSecretShared
}

// RpcEventType represents any RPC event
type RpcEventType interface {
	GetMethod() RpcMethod
}

// PingRpcEvent represents a ping RPC event
type PingRpcEvent struct {
	Method RpcMethod `json:"method"`
	Result string    `json:"result,omitempty"`
}

func (e *PingRpcEvent) GetMethod() RpcMethod {
	return RpcMethodPing
}

// GetActiveOrdersRpcEvent represents a getActiveOrders RPC event
type GetActiveOrdersRpcEvent struct {
	Method RpcMethod   `json:"method"`
	Result interface{} `json:"result"`
}

func (e *GetActiveOrdersRpcEvent) GetMethod() RpcMethod {
	return RpcMethodGetActiveOrders
}

// GetSecretsRpcEvent represents a getSecrets RPC event
type GetSecretsRpcEvent struct {
	Method RpcMethod   `json:"method"`
	Result interface{} `json:"result"`
}

func (e *GetSecretsRpcEvent) GetMethod() RpcMethod {
	return RpcMethodGetSecrets
}

// GetAllowedMethodsRpcEvent represents a getAllowedMethods RPC event
type GetAllowedMethodsRpcEvent struct {
	Method RpcMethod `json:"method"`
	Result []string  `json:"result"`
}

func (e *GetAllowedMethodsRpcEvent) GetMethod() RpcMethod {
	return RpcMethodGetAllowedMethods
}

// Callback types
type OnOrderCb func(data OrderEventType)
type OnOrderCreatedCb func(data *OrderCreatedEvent)
type OnOrderInvalidCb func(data *OrderInvalidEvent)
type OnOrderBalanceChangeCb func(data *OrderBalanceChangeEvent)
type OnOrderAllowanceChangeCb func(data *OrderAllowanceChangeEvent)
type OnOrderFilledCb func(data *OrderFilledEvent)
type OnOrderFilledPartiallyCb func(data *OrderFilledPartiallyEvent)
type OnOrderCancelledCb func(data *OrderCancelledEvent)
type OnOrderSecretSharedCb func(data *OrderSecretSharedEvent)
type OnPongCb func(result string)
type OnGetActiveOrdersCb func(result interface{})
type OnGetSecretsCb func(result interface{})
type OnGetAllowedMethodsCb func(result []string)
type OnMessageCb func(data interface{})
type OnErrorCb func(err error)
type OnCloseCb func()
type OnOpenCb func()
