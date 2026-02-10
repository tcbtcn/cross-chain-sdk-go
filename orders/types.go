package orders

import (
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/dawitel/cross-chain-sdk-go/chains"
	"github.com/dawitel/cross-chain-sdk-go/crypto/eip712"
	"github.com/dawitel/cross-chain-sdk-go/domains/addresses"
	"github.com/dawitel/cross-chain-sdk-go/domains/auction"
	"github.com/dawitel/cross-chain-sdk-go/domains/hashlock"
	"github.com/dawitel/cross-chain-sdk-go/domains/timelocks"
)

type LimitOrderV4Struct struct {
	Maker         string `json:"maker"`
	MakerAsset    string `json:"makerAsset"`
	TakerAsset    string `json:"takerAsset"`
	MakingAmount  string `json:"makingAmount"`
	TakingAmount  string `json:"takingAmount"`
	Receiver      string `json:"receiver"`
	AllowedSender string `json:"allowedSender"`
	MakerTraits   string `json:"makerTraits"`
	Salt          string `json:"salt"`
	Expiration    string `json:"expiration"`
	Nonce         string `json:"nonce"`
}

type EvmCrossChainOrder struct {
	Maker                *addresses.EvmAddress
	MakerAsset           *addresses.EvmAddress
	TakerAsset           addresses.AddressLike
	MakingAmount         *big.Int
	TakingAmount         *big.Int
	Receiver             addresses.AddressLike
	Deadline             *big.Int
	AuctionStartTime     *big.Int
	AuctionEndTime       *big.Int
	Nonce                *big.Int
	Salt                 *big.Int
	HashLock             *hashlock.HashLock
	TimeLocks            *timelocks.TimeLocks
	SrcSafetyDeposit     *big.Int
	DstSafetyDeposit     *big.Int
	DstChainID           chains.SupportedChain
	Extension            string
	MultipleFillsAllowed bool
}

// validateRequiredFields checks that all required fields for order operations are non-nil
func (o *EvmCrossChainOrder) validateRequiredFields() error {
	if o == nil {
		return fmt.Errorf("order is nil")
	}
	if o.Maker == nil {
		return fmt.Errorf("maker address is required")
	}
	if o.MakerAsset == nil {
		return fmt.Errorf("makerAsset address is required")
	}
	if o.TakerAsset == nil {
		return fmt.Errorf("takerAsset address is required")
	}
	if o.Receiver == nil {
		return fmt.Errorf("receiver address is required")
	}
	if o.MakingAmount == nil {
		return fmt.Errorf("makingAmount is required")
	}
	if o.TakingAmount == nil {
		return fmt.Errorf("takingAmount is required")
	}
	if o.Deadline == nil {
		return fmt.Errorf("deadline is required")
	}
	if o.Nonce == nil {
		return fmt.Errorf("nonce is required")
	}
	if o.Salt == nil {
		return fmt.Errorf("salt is required")
	}
	if o.HashLock == nil {
		return fmt.Errorf("hashLock is required")
	}
	return nil
}

func (o *EvmCrossChainOrder) GetOrderHash(srcChainID chains.SupportedChain) (string, error) {
	if err := o.validateRequiredFields(); err != nil {
		return "", fmt.Errorf("invalid order: %w", err)
	}
	typedData, err := o.GetTypedData(srcChainID)
	if err != nil {
		return "", fmt.Errorf("failed to get typed data: %w", err)
	}
	typedDataMap, ok := typedData.(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid typed data format")
	}

	// Convert types from []map[string]string to []interface{}
	typesMap := typedDataMap["types"].(map[string]interface{})
	convertedTypes := make(map[string]interface{})
	for typeName, typeValue := range typesMap {
		if typeSlice, ok := typeValue.([]map[string]string); ok {
			convertedSlice := make([]interface{}, len(typeSlice))
			for i, v := range typeSlice {
				convertedSlice[i] = v
			}
			convertedTypes[typeName] = convertedSlice
		} else {
			convertedTypes[typeName] = typeValue
		}
	}

	eip712Data := eip712.TypedData{
		Types:       convertedTypes,
		PrimaryType: typedDataMap["primaryType"].(string),
		Domain:      typedDataMap["domain"].(map[string]interface{}),
		Message:     typedDataMap["message"].(map[string]interface{}),
	}

	hash, err := eip712.HashTypedData(eip712Data)
	if err != nil {
		return "", fmt.Errorf("failed to hash typed data: %w", err)
	}
	return "0x" + hex.EncodeToString(hash), nil
}

func (o *EvmCrossChainOrder) GetTypedData(srcChainID chains.SupportedChain) (interface{}, error) {
	if err := o.validateRequiredFields(); err != nil {
		return nil, fmt.Errorf("invalid order: %w", err)
	}
	return map[string]interface{}{
		"types": map[string]interface{}{
			"EIP712Domain": []map[string]string{
				{"name": "name", "type": "string"},
				{"name": "version", "type": "string"},
				{"name": "chainId", "type": "uint256"},
				{"name": "verifyingContract", "type": "address"},
			},
			"Order": []map[string]string{
				{"name": "salt", "type": "uint256"},
				{"name": "maker", "type": "address"},
				{"name": "receiver", "type": "address"},
				{"name": "makerAsset", "type": "address"},
				{"name": "takerAsset", "type": "address"},
				{"name": "makingAmount", "type": "uint256"},
				{"name": "takingAmount", "type": "uint256"},
				{"name": "makerTraits", "type": "uint256"},
			},
		},
		"primaryType": "Order",
		"domain": map[string]interface{}{
			"name":              "1inch Limit Order Protocol",
			"version":           "4",
			"chainId":           int(srcChainID),
			"verifyingContract": "0x0000000000000000000000000000000000000000",
		},
		"message": map[string]interface{}{
			"salt":         o.Salt.String(),
			"maker":        o.Maker.ToString(),
			"receiver":     o.Receiver.ToString(),
			"makerAsset":   o.MakerAsset.ToString(),
			"takerAsset":   o.TakerAsset.ToString(),
			"makingAmount": o.MakingAmount.String(),
			"takingAmount": o.TakingAmount.String(),
			"makerTraits":  "0",
		},
	}, nil
}

func (o *EvmCrossChainOrder) Build() (LimitOrderV4Struct, error) {
	if err := o.validateRequiredFields(); err != nil {
		return LimitOrderV4Struct{}, fmt.Errorf("invalid order: %w", err)
	}
	return LimitOrderV4Struct{
		Maker:         o.Maker.ToString(),
		MakerAsset:    o.MakerAsset.ToString(),
		TakerAsset:    o.TakerAsset.ToString(),
		MakingAmount:  o.MakingAmount.String(),
		TakingAmount:  o.TakingAmount.String(),
		Receiver:      o.Receiver.ToString(),
		AllowedSender: "0x0000000000000000000000000000000000000000",
		MakerTraits:   "0",
		Salt:          o.Salt.String(),
		Expiration:    o.Deadline.String(),
		Nonce:         o.Nonce.String(),
	}, nil
}

type SolanaCrossChainOrder struct {
	OrderHash            []byte
	HashLock             *hashlock.HashLock
	Auction              *auction.AuctionDetails
	MultipleFillsAllowed bool
	// Additional fields for proper order structure
	SrcToken         *addresses.SolanaAddress
	DstToken         *addresses.EvmAddress
	Maker            *addresses.SolanaAddress
	Receiver         *addresses.EvmAddress
	SrcAmount        *big.Int
	MinDstAmount     *big.Int
	SrcSafetyDeposit *big.Int
	DstSafetyDeposit *big.Int
	TimeLocks        *timelocks.TimeLocks
	DstChainID       chains.SupportedChain
	Salt             *big.Int
	Source           string
	SrcAssetIsNative bool
}

func (o *SolanaCrossChainOrder) GetOrderHash() string {
	return "0x" + hex.EncodeToString(o.OrderHash)
}

func (o *SolanaCrossChainOrder) ToJSON() map[string]interface{} {
	auctionJSON := map[string]interface{}{}
	if o.Auction != nil {
		auctionJSON["startTime"] = o.Auction.StartTime.String()
		auctionJSON["duration"] = o.Auction.Duration.String()
		auctionJSON["initialRateBump"] = o.Auction.InitialRateBump
		points := make([]map[string]interface{}, len(o.Auction.Points))
		for i, p := range o.Auction.Points {
			points[i] = map[string]interface{}{
				"toTokenAmount": p.ToTokenAmount,
				"delay":         p.Delay,
			}
		}
		auctionJSON["points"] = points
	}

	orderInfo := map[string]interface{}{}
	if o.SrcToken != nil {
		orderInfo["srcToken"] = o.SrcToken.ToString()
	}
	if o.DstToken != nil {
		orderInfo["dstToken"] = o.DstToken.ToString()
	}
	if o.Maker != nil {
		orderInfo["maker"] = o.Maker.ToString()
	}
	if o.Receiver != nil {
		orderInfo["receiver"] = o.Receiver.ToString()
	}
	if o.SrcAmount != nil {
		orderInfo["srcAmount"] = o.SrcAmount.String()
	}
	if o.MinDstAmount != nil {
		orderInfo["minDstAmount"] = o.MinDstAmount.String()
	}

	escrowParams := map[string]interface{}{}
	if o.HashLock != nil {
		escrowParams["hashLock"] = o.HashLock.ToString()
	}
	escrowParams["srcChainId"] = int(chains.Solana)
	escrowParams["dstChainId"] = int(o.DstChainID)
	if o.SrcSafetyDeposit != nil {
		escrowParams["srcSafetyDeposit"] = o.SrcSafetyDeposit.String()
	}
	if o.DstSafetyDeposit != nil {
		escrowParams["dstSafetyDeposit"] = o.DstSafetyDeposit.String()
	}
	if o.TimeLocks != nil {
		timeLocksVal := o.TimeLocks.Build()
		escrowParams["timeLocks"] = timeLocksVal.String()
	}

	extra := map[string]interface{}{
		"srcAssetIsNative":    o.SrcAssetIsNative,
		"allowMultipleFills":   o.MultipleFillsAllowed,
		"orderExpirationDelay": "12", // Default value
		"resolverCancellationConfig": map[string]interface{}{
			"maxCancellationPremium":      "0",
			"cancellationAuctionDuration": 0,
		},
	}
	if o.Source != "" {
		extra["source"] = o.Source
	} else {
		extra["source"] = "sdk"
	}
	if o.Salt != nil {
		extra["salt"] = o.Salt.String()
	}

	return map[string]interface{}{
		"orderInfo":    orderInfo,
		"escrowParams": escrowParams,
		"details": map[string]interface{}{
			"auction": auctionJSON,
		},
		"extra": extra,
	}
}
