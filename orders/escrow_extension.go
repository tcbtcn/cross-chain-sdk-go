package orders

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common/math"
	"github.com/tcbtcn/cross-chain-sdk-go/chains"
	"github.com/tcbtcn/cross-chain-sdk-go/domains/addresses"
	"github.com/tcbtcn/cross-chain-sdk-go/domains/auction"
	"github.com/tcbtcn/cross-chain-sdk-go/domains/hashlock"
	"github.com/tcbtcn/cross-chain-sdk-go/domains/timelocks"
)

const ZX = "0x"

type EscrowExtension struct {
	Address             *addresses.EvmAddress
	AuctionDetails      *auction.AuctionDetails
	PostInteractionData *SettlementPostInteractionData
	MakerPermit         []byte
	HashLock            *hashlock.HashLock
	DstChainID          chains.SupportedChain
	DstToken            addresses.AddressLike
	SrcSafetyDeposit    *big.Int
	DstSafetyDeposit    *big.Int
	TimeLocks           *timelocks.TimeLocks
	DstAddressFirstPart *addresses.AddressComplement
}

func NewEscrowExtension(
	address *addresses.EvmAddress,
	auctionDetails *auction.AuctionDetails,
	postInteractionData *SettlementPostInteractionData,
	makerPermit []byte,
	hashLock *hashlock.HashLock,
	dstChainID chains.SupportedChain,
	dstToken addresses.AddressLike,
	srcSafetyDeposit *big.Int,
	dstSafetyDeposit *big.Int,
	timeLocks *timelocks.TimeLocks,
	dstAddressFirstPart *addresses.AddressComplement,
) *EscrowExtension {
	return &EscrowExtension{
		Address:             address,
		AuctionDetails:      auctionDetails,
		PostInteractionData: postInteractionData,
		MakerPermit:         makerPermit,
		HashLock:            hashLock,
		DstChainID:          dstChainID,
		DstToken:            dstToken,
		SrcSafetyDeposit:    srcSafetyDeposit,
		DstSafetyDeposit:    dstSafetyDeposit,
		TimeLocks:           timeLocks,
		DstAddressFirstPart: dstAddressFirstPart,
	}
}

func (ee *EscrowExtension) Build() (string, error) {
	baseExt, err := ee.buildBaseExtension()
	if err != nil {
		return "", fmt.Errorf("failed to build base extension: %w", err)
	}

	extraData, err := EncodeExtraEscrowData(
		ee.HashLock,
		int(ee.DstChainID),
		ee.DstToken,
		ee.SrcSafetyDeposit,
		ee.DstSafetyDeposit,
		ee.TimeLocks,
	)
	if err != nil {
		return "", fmt.Errorf("failed to encode extra escrow data: %w", err)
	}

	postInteractionWithExtra := baseExt.PostInteraction + extraData

	customData := ee.buildCustomData()

	extension, err := ee.encodeExtension(
		baseExt.Address,
		baseExt.AuctionDetails,
		postInteractionWithExtra,
		customData,
	)
	if err != nil {
		return "", fmt.Errorf("failed to encode extension: %w", err)
	}

	return extension, nil
}

type BaseExtension struct {
	Address         []byte
	AuctionDetails  []byte
	PostInteraction string
	CustomData      string
}

func (ee *EscrowExtension) buildBaseExtension() (*BaseExtension, error) {
	if ee.Address == nil {
		return nil, fmt.Errorf("address is required")
	}
	addressBytes := ee.Address.ToBuffer()

	auctionDetailsBytes, err := ee.AuctionDetails.EncodeForEvm()
	if err != nil {
		return nil, fmt.Errorf("failed to encode auction details: %w", err)
	}

	postInteraction, err := ee.PostInteractionData.EncodeForExtension(ee.Address)
	if err != nil {
		return nil, fmt.Errorf("failed to encode post interaction for extension: %w", err)
	}

	customData := ee.buildCustomData()

	return &BaseExtension{
		Address:         addressBytes,
		AuctionDetails:  auctionDetailsBytes,
		PostInteraction: postInteraction,
		CustomData:      customData,
	}, nil
}

func (ee *EscrowExtension) buildCustomData() string {
	if ee.DstAddressFirstPart == nil || ee.DstAddressFirstPart.IsZero() {
		return ZX
	}
	return ee.DstAddressFirstPart.AsHex()
}

func (ee *EscrowExtension) encodeExtension(
	address []byte,
	auctionDetails []byte,
	postInteraction string,
	customData string,
) (string, error) {
	postInteractionBytes, err := hex.DecodeString(strings.TrimPrefix(postInteraction, "0x"))
	if err != nil {
		return "", fmt.Errorf("failed to decode postInteraction: %w", err)
	}

	var customDataBytes []byte
	if customData != ZX && customData != "" {
		customDataBytes, err = hex.DecodeString(strings.TrimPrefix(customData, "0x"))
		if err != nil {
			return "", fmt.Errorf("failed to decode customData: %w", err)
		}
	}

	addressOffset := uint32(32)
	auctionDetailsOffset := addressOffset + uint32(len(address))
	postInteractionOffset := auctionDetailsOffset + uint32(len(auctionDetails))
	customDataOffset := postInteractionOffset + uint32(len(postInteractionBytes))

	offsets := make([]byte, 32)
	copy(offsets[0:8], math.U256Bytes(big.NewInt(int64(addressOffset))))
	copy(offsets[8:16], math.U256Bytes(big.NewInt(int64(auctionDetailsOffset))))
	copy(offsets[16:24], math.U256Bytes(big.NewInt(int64(postInteractionOffset))))
	copy(offsets[24:32], math.U256Bytes(big.NewInt(int64(customDataOffset))))

	var result []byte
	result = append(result, offsets...)
	result = append(result, address...)
	result = append(result, auctionDetails...)
	result = append(result, postInteractionBytes...)
	if len(customDataBytes) > 0 {
		result = append(result, customDataBytes...)
	}

	return "0x" + hex.EncodeToString(result), nil
}

func (ee *EscrowExtension) Encode() (string, error) {
	return ee.Build()
}
