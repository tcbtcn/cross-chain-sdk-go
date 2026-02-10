package orders

import (
	"github.com/dawitel/cross-chain-sdk-go/chains"
	"github.com/dawitel/cross-chain-sdk-go/orders"
)

type OrdersApiConfig struct {
	URL     string
	AuthKey string
}

type ActiveOrdersRequestParams struct {
	Page       int
	Limit      int
	SrcChainID *chains.SupportedChain
	DstChainID *chains.SupportedChain
}

type FillInfo struct {
	TxHash string `json:"txHash"`
}

type ActiveOrder struct {
	QuoteID              string                `json:"quoteId"`
	OrderHash            string                `json:"orderHash"`
	Deadline             string                `json:"deadline"`
	AuctionStartDate     string                `json:"auctionStartDate"`
	AuctionEndDate       string                `json:"auctionEndDate"`
	DstChainID           chains.SupportedChain `json:"dstChainId"`
	RemainingMakerAmount string                `json:"remainingMakerAmount"`
	SecretHashes         []string              `json:"secretHashes,omitempty"`
	Fills                []FillInfo            `json:"fills"`
}

type ActiveOrdersResponse struct {
	Items      []ActiveOrder `json:"items"`
	TotalCount int           `json:"totalCount"`
	Page       int           `json:"page"`
	Limit      int           `json:"limit"`
}

type OrderStatusParams struct {
	OrderHash string
}

type ValidationStatus string

const (
	ValidationStatusValid                       ValidationStatus = "valid"
	ValidationStatusOrderPredicateReturnedFalse ValidationStatus = "order-predicate-returned-false"
	ValidationStatusNotEnoughBalance            ValidationStatus = "not-enough-balance"
	ValidationStatusNotEnoughAllowance          ValidationStatus = "not-enough-allowance"
	ValidationStatusInvalidPermitSignature      ValidationStatus = "invalid-permit-signature"
	ValidationStatusInvalidPermitSpender        ValidationStatus = "invalid-permit-spender"
	ValidationStatusInvalidPermitSigner         ValidationStatus = "invalid-permit-signer"
	ValidationStatusInvalidSignature            ValidationStatus = "invalid-signature"
	ValidationStatusFailedToParsePermitDetails  ValidationStatus = "failed-to-parse-permit-details"
	ValidationStatusUnknownPermitVersion        ValidationStatus = "unknown-permit-version"
	ValidationStatusWrongEpochManagerAndBit     ValidationStatus = "wrong-epoch-manager-and-bit-invalidator"
	ValidationStatusFailedToDecodeRemaining     ValidationStatus = "failed-to-decode-remaining"
	ValidationStatusUnknownFailure              ValidationStatus = "unknown-failure"
)

type FillStatus string

const (
	FillStatusPending   FillStatus = "pending"
	FillStatusExecuted  FillStatus = "executed"
	FillStatusRefunding FillStatus = "refunding"
	FillStatusRefunded  FillStatus = "refunded"
)

type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusExecuted  OrderStatus = "executed"
	OrderStatusExpired   OrderStatus = "expired"
	OrderStatusCancelled OrderStatus = "cancelled"
	OrderStatusRefunding OrderStatus = "refunding"
	OrderStatusRefunded  OrderStatus = "refunded"
)

type Fill struct {
	Status                   FillStatus        `json:"status"`
	TxHash                   string            `json:"txHash"`
	FilledMakerAmount        string            `json:"filledMakerAmount"`
	FilledAuctionTakerAmount string            `json:"filledAuctionTakerAmount"`
	EscrowEvents             []EscrowEventData `json:"escrowEvents"`
}

type EscrowEventSide string

const (
	EscrowEventSideSrc EscrowEventSide = "src"
	EscrowEventSideDst EscrowEventSide = "dst"
)

type EscrowEventAction string

const (
	EscrowEventActionSrcEscrowCreated EscrowEventAction = "src_escrow_created"
	EscrowEventActionDstEscrowCreated EscrowEventAction = "dst_escrow_created"
	EscrowEventActionWithdrawn        EscrowEventAction = "withdrawn"
	EscrowEventActionFundsRescued     EscrowEventAction = "funds_rescued"
	EscrowEventActionEscrowCancelled  EscrowEventAction = "escrow_cancelled"
)

type EscrowEventData struct {
	TransactionHash string            `json:"transactionHash"`
	Escrow          string            `json:"escrow"`
	Side            EscrowEventSide   `json:"side"`
	Action          EscrowEventAction `json:"action"`
	BlockTimestamp  int64             `json:"blockTimestamp"`
}

type OrderStatusResponse struct {
	OrderHash        string                     `json:"orderHash"`
	Status           OrderStatus                `json:"status"`
	ValidationStatus ValidationStatus           `json:"validationStatus"`
	Fills            []Fill                     `json:"fills"`
	SrcChainID       chains.SupportedChain      `json:"srcChainId"`
	DstChainID       chains.SupportedChain      `json:"dstChainId"`
	Order            *orders.LimitOrderV4Struct `json:"order,omitempty"`
	Extension        string                     `json:"extension,omitempty"`
}

type OrdersByMakerParams struct {
	Address string
	Page    int
	Limit   int
}

type OrderFillsByMakerOutput struct {
	OrderHash string `json:"orderHash"`
	Fills     []Fill `json:"fills"`
}

type OrdersByMakerResponse struct {
	Items      []OrderFillsByMakerOutput `json:"items"`
	TotalCount int                       `json:"totalCount"`
	Page       int                       `json:"page"`
	Limit      int                       `json:"limit"`
}

type ReadyToAcceptSecretFill struct {
	Idx int `json:"idx"`
}

type ReadyToAcceptSecretFills struct {
	Fills []ReadyToAcceptSecretFill `json:"fills"`
}

type ChainImmutables struct {
	ChainID chains.SupportedChain `json:"chainId"`
	Factory string                `json:"factory"`
}

type PublicSecret struct {
	SecretHash string `json:"secretHash"`
	Secret     string `json:"secret"`
}

type PublishedSecretsResponse struct {
	Secrets []PublicSecret `json:"secrets"`
}

type CancellableOrdersResponse struct {
	Items      []interface{} `json:"items"`
	TotalCount int           `json:"totalCount"`
	Page       int           `json:"page"`
	Limit      int           `json:"limit"`
}

type ReadyToExecutePublicAction struct {
	OrderHash string `json:"orderHash"`
	Action    string `json:"action"`
}

type ReadyToExecutePublicActions struct {
	Actions []ReadyToExecutePublicAction `json:"actions"`
}
