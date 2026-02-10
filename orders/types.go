package orders

import (
	"encoding/hex"
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

func (o *EvmCrossChainOrder) GetOrderHash(srcChainID chains.SupportedChain) string {
	if o.Maker == nil || o.Receiver == nil || o.MakerAsset == nil || o.TakerAsset == nil {
		return ""
	}
	typedData := o.GetTypedData(srcChainID)
	typedDataMap, ok := typedData.(map[string]interface{})
	if !ok {
		return ""
	}

	eip712Data := eip712.TypedData{
		Types:       typedDataMap["types"].(map[string]interface{}),
		PrimaryType: typedDataMap["primaryType"].(string),
		Domain:      typedDataMap["domain"].(map[string]interface{}),
		Message:     typedDataMap["message"].(map[string]interface{}),
	}

	hash, err := eip712.HashTypedData(eip712Data)
	if err != nil {
		return ""
	}
	return "0x" + hex.EncodeToString(hash)
}

func (o *EvmCrossChainOrder) GetTypedData(srcChainID chains.SupportedChain) interface{} {
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
	}
}

func (o *EvmCrossChainOrder) Build() LimitOrderV4Struct {
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
	}
}

type SolanaCrossChainOrder struct {
	OrderHash            []byte
	HashLock             *hashlock.HashLock
	Auction              *auction.AuctionDetails
	MultipleFillsAllowed bool
}

func (o *SolanaCrossChainOrder) GetOrderHash() string {
	return "0x" + hex.EncodeToString(o.OrderHash)
}

func (o *SolanaCrossChainOrder) ToJSON() map[string]interface{} {
	return map[string]interface{}{
		"orderHash": o.GetOrderHash(),
		"hashLock":  o.HashLock.ToString(),
	}
}
