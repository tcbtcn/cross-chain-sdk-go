package wsapi

// OrderEvents is a list of all order event types
var OrderEvents = []EventType{
	EventTypeOrderCreated,
	EventTypeOrderInvalid,
	EventTypeOrderBalanceChange,
	EventTypeOrderAllowanceChange,
	EventTypeOrderFilled,
	EventTypeOrderFilledPartially,
	EventTypeOrderCancelled,
	EventTypeOrderSecretShared,
}
