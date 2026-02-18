package orders

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/tcbtcn/cross-chain-sdk-go/domains/addresses"
)

type SettlementPostInteractionData struct {
	Whitelist          []AuctionWhitelistItem
	ResolvingStartTime *big.Int
}

func NewSettlementPostInteractionData(whitelist []AuctionWhitelistItem, resolvingStartTime *big.Int) *SettlementPostInteractionData {
	return &SettlementPostInteractionData{
		Whitelist:          whitelist,
		ResolvingStartTime: resolvingStartTime,
	}
}

func (spid *SettlementPostInteractionData) Encode() (string, error) {
	var result []byte

	whitelistData, err := encodeWhitelist(spid.Whitelist)
	if err != nil {
		return "", fmt.Errorf("failed to encode whitelist: %w", err)
	}

	resolvingStartTimeBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(resolvingStartTimeBytes, uint32(spid.ResolvingStartTime.Uint64()))
	result = append(result, resolvingStartTimeBytes...)

	result = append(result, whitelistData...)

	bitmap := buildExtensionsBitmap(len(spid.Whitelist), 0)
	result = append(result, bitmap)

	return "0x" + hex.EncodeToString(result), nil
}

func encodeWhitelist(whitelist []AuctionWhitelistItem) ([]byte, error) {
	var result []byte

	for _, item := range whitelist {
		if item.Address == nil {
			return nil, fmt.Errorf("whitelist item address is nil")
		}
		addrBytes := item.Address.ToBuffer()
		if len(addrBytes) != 20 {
			return nil, fmt.Errorf("invalid address length: %d", len(addrBytes))
		}

		result = append(result, addrBytes[:10]...)

		if item.AllowFrom == nil {
			return nil, fmt.Errorf("whitelist item allowFrom is nil")
		}
		allowFromBytes := make([]byte, 2)
		allowFromUint64 := item.AllowFrom.Uint64()
		if allowFromUint64 > 0xffff {
			return nil, fmt.Errorf("allowFrom exceeds uint16 max")
		}
		binary.BigEndian.PutUint16(allowFromBytes, uint16(allowFromUint64))
		result = append(result, allowFromBytes...)
	}

	return result, nil
}

func buildExtensionsBitmap(resolverCount int, feeType int) byte {
	if resolverCount > 31 {
		panic("resolver count exceeds 31")
	}
	if feeType > 7 {
		panic("fee type exceeds 7")
	}

	var bitmap byte
	bitmap |= byte(feeType)
	bitmap |= byte(resolverCount << 3)

	return bitmap
}

func (spid *SettlementPostInteractionData) EncodeForExtension(settlementAddress *addresses.EvmAddress) (string, error) {
	if settlementAddress == nil {
		return "", fmt.Errorf("settlement address is required")
	}

	postInteractionData, err := spid.Encode()
	if err != nil {
		return "", err
	}

	settlementAddrHex := settlementAddress.ToHex()
	postInteractionDataTrimmed := strings.TrimPrefix(postInteractionData, "0x")

	return settlementAddrHex + postInteractionDataTrimmed, nil
}

func DecodeSettlementPostInteractionData(data string) (*SettlementPostInteractionData, error) {
	data = strings.TrimPrefix(data, "0x")
	if len(data) < 2 {
		return nil, fmt.Errorf("data too short")
	}

	bytes, err := hex.DecodeString(data)
	if err != nil {
		return nil, fmt.Errorf("failed to decode hex: %w", err)
	}

	if len(bytes) < 4 {
		return nil, fmt.Errorf("data too short for resolvingStartTime")
	}

	resolvingStartTime := big.NewInt(int64(binary.BigEndian.Uint32(bytes[0:4])))

	bitmap := bytes[len(bytes)-1]
	resolverCount := int((bitmap >> 3) & 0x1f)

	whitelistStart := 4
	whitelistEnd := len(bytes) - 1
	whitelistBytes := bytes[whitelistStart:whitelistEnd]

	whitelist := make([]AuctionWhitelistItem, resolverCount)
	for i := 0; i < resolverCount; i++ {
		offset := i * 12
		if offset+12 > len(whitelistBytes) {
			return nil, fmt.Errorf("whitelist data too short")
		}

		addrBytes := make([]byte, 20)
		copy(addrBytes[10:], whitelistBytes[offset:offset+10])
		addr := common.BytesToAddress(addrBytes)
		evmAddr := addresses.NewEvmAddress(addr.Hex())

		allowFromBytes := whitelistBytes[offset+10 : offset+12]
		allowFrom := big.NewInt(int64(binary.BigEndian.Uint16(allowFromBytes)))

		whitelist[i] = AuctionWhitelistItem{
			Address:   evmAddr,
			AllowFrom: allowFrom,
		}
	}

	return &SettlementPostInteractionData{
		Whitelist:          whitelist,
		ResolvingStartTime: resolvingStartTime,
	}, nil
}
