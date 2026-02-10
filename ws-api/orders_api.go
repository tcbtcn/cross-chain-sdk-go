package wsapi

import (
	"encoding/json"
)

// ActiveOrdersWebSocketApi handles order-related WebSocket events
type ActiveOrdersWebSocketApi struct {
	provider WsProviderConnector
}

// NewActiveOrdersWebSocketApi creates a new ActiveOrdersWebSocketApi
func NewActiveOrdersWebSocketApi(provider WsProviderConnector) *ActiveOrdersWebSocketApi {
	return &ActiveOrdersWebSocketApi{
		provider: provider,
	}
}

// OnOrder subscribes to all order events
func (a *ActiveOrdersWebSocketApi) OnOrder(cb OnOrderCb) {
	a.provider.OnMessage(func(data interface{}) {
		eventData, ok := a.parseOrderEvent(data)
		if ok {
			cb(eventData)
		}
	})
}

// OnOrderCreated subscribes to order_created events
func (a *ActiveOrdersWebSocketApi) OnOrderCreated(cb OnOrderCreatedCb) {
	a.provider.OnMessage(func(data interface{}) {
		eventData, ok := a.parseOrderEvent(data)
		if ok && eventData.GetEvent() == EventTypeOrderCreated {
			if createdEvent, ok := eventData.(*OrderCreatedEvent); ok {
				cb(createdEvent)
			}
		}
	})
}

// OnOrderInvalid subscribes to order_invalid events
func (a *ActiveOrdersWebSocketApi) OnOrderInvalid(cb OnOrderInvalidCb) {
	a.provider.OnMessage(func(data interface{}) {
		eventData, ok := a.parseOrderEvent(data)
		if ok && eventData.GetEvent() == EventTypeOrderInvalid {
			if invalidEvent, ok := eventData.(*OrderInvalidEvent); ok {
				cb(invalidEvent)
			}
		}
	})
}

// OnOrderBalanceChange subscribes to order_balance_change events
func (a *ActiveOrdersWebSocketApi) OnOrderBalanceChange(cb OnOrderBalanceChangeCb) {
	a.provider.OnMessage(func(data interface{}) {
		eventData, ok := a.parseOrderEvent(data)
		if ok && eventData.GetEvent() == EventTypeOrderBalanceChange {
			if balanceEvent, ok := eventData.(*OrderBalanceChangeEvent); ok {
				cb(balanceEvent)
			}
		}
	})
}

// OnOrderAllowanceChange subscribes to order_allowance_change events
func (a *ActiveOrdersWebSocketApi) OnOrderAllowanceChange(cb OnOrderAllowanceChangeCb) {
	a.provider.OnMessage(func(data interface{}) {
		eventData, ok := a.parseOrderEvent(data)
		if ok && eventData.GetEvent() == EventTypeOrderAllowanceChange {
			if allowanceEvent, ok := eventData.(*OrderAllowanceChangeEvent); ok {
				cb(allowanceEvent)
			}
		}
	})
}

// OnOrderFilled subscribes to order_filled events
func (a *ActiveOrdersWebSocketApi) OnOrderFilled(cb OnOrderFilledCb) {
	a.provider.OnMessage(func(data interface{}) {
		eventData, ok := a.parseOrderEvent(data)
		if ok && eventData.GetEvent() == EventTypeOrderFilled {
			if filledEvent, ok := eventData.(*OrderFilledEvent); ok {
				cb(filledEvent)
			}
		}
	})
}

// OnOrderFilledPartially subscribes to order_filled_partially events
func (a *ActiveOrdersWebSocketApi) OnOrderFilledPartially(cb OnOrderFilledPartiallyCb) {
	a.provider.OnMessage(func(data interface{}) {
		eventData, ok := a.parseOrderEvent(data)
		if ok && eventData.GetEvent() == EventTypeOrderFilledPartially {
			if partialEvent, ok := eventData.(*OrderFilledPartiallyEvent); ok {
				cb(partialEvent)
			}
		}
	})
}

// OnOrderCancelled subscribes to order_cancelled events
func (a *ActiveOrdersWebSocketApi) OnOrderCancelled(cb OnOrderCancelledCb) {
	a.provider.OnMessage(func(data interface{}) {
		eventData, ok := a.parseOrderEvent(data)
		if ok && eventData.GetEvent() == EventTypeOrderCancelled {
			if cancelledEvent, ok := eventData.(*OrderCancelledEvent); ok {
				cb(cancelledEvent)
			}
		}
	})
}

// OnOrderSecretShared subscribes to secret_shared events
func (a *ActiveOrdersWebSocketApi) OnOrderSecretShared(cb OnOrderSecretSharedCb) {
	a.provider.OnMessage(func(data interface{}) {
		eventData, ok := a.parseOrderEvent(data)
		if ok && eventData.GetEvent() == EventTypeOrderSecretShared {
			if secretEvent, ok := eventData.(*OrderSecretSharedEvent); ok {
				cb(secretEvent)
			}
		}
	})
}

func (a *ActiveOrdersWebSocketApi) parseOrderEvent(data interface{}) (OrderEventType, bool) {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return nil, false
	}

	var eventMap map[string]interface{}
	if err := json.Unmarshal(dataBytes, &eventMap); err != nil {
		return nil, false
	}

	eventTypeStr, ok := eventMap["event"].(string)
	if !ok {
		return nil, false
	}

	eventType := EventType(eventTypeStr)

	switch eventType {
	case EventTypeOrderCreated:
		var event OrderCreatedEvent
		if err := json.Unmarshal(dataBytes, &event); err == nil {
			return &event, true
		}
	case EventTypeOrderInvalid:
		var event OrderInvalidEvent
		if err := json.Unmarshal(dataBytes, &event); err == nil {
			return &event, true
		}
	case EventTypeOrderBalanceChange:
		var event OrderBalanceChangeEvent
		if err := json.Unmarshal(dataBytes, &event); err == nil {
			return &event, true
		}
	case EventTypeOrderAllowanceChange:
		var event OrderAllowanceChangeEvent
		if err := json.Unmarshal(dataBytes, &event); err == nil {
			return &event, true
		}
	case EventTypeOrderFilled:
		var event OrderFilledEvent
		if err := json.Unmarshal(dataBytes, &event); err == nil {
			return &event, true
		}
	case EventTypeOrderFilledPartially:
		var event OrderFilledPartiallyEvent
		if err := json.Unmarshal(dataBytes, &event); err == nil {
			return &event, true
		}
	case EventTypeOrderCancelled:
		var event OrderCancelledEvent
		if err := json.Unmarshal(dataBytes, &event); err == nil {
			return &event, true
		}
	case EventTypeOrderSecretShared:
		var event OrderSecretSharedEvent
		if err := json.Unmarshal(dataBytes, &event); err == nil {
			return &event, true
		}
	}

	return nil, false
}
