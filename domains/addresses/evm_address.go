package addresses

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

const (
	ZeroAddress    = "0x0000000000000000000000000000000000000000"
	NativeCurrency = "0xeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee"
)

var (
	EvmZero   = NewEvmAddress(ZeroAddress)
	EvmNative = NewEvmAddress(NativeCurrency)
)

type EvmAddress struct {
	inner common.Address
}

func NewEvmAddress(address string) *EvmAddress {
	addr := common.HexToAddress(address)
	return &EvmAddress{inner: addr}
}

func EvmAddressFromString(address string) (*EvmAddress, error) {
	if !common.IsHexAddress(address) {
		return nil, fmt.Errorf("invalid address: %s", address)
	}
	return NewEvmAddress(address), nil
}

func EvmAddressFromBigInt(val *big.Int) *EvmAddress {
	addr := common.BigToAddress(val)
	return &EvmAddress{inner: addr}
}

func EvmAddressFromBuffer(buf []byte) (*EvmAddress, error) {
	if len(buf) < 20 {
		return nil, fmt.Errorf("buffer too short for address")
	}
	var addr common.Address
	copy(addr[:], buf[len(buf)-20:])
	return &EvmAddress{inner: addr}, nil
}

func (a *EvmAddress) NativeAsZero() AddressLike {
	if a.IsNative() {
		return EvmZero
	}
	return a
}

func (a *EvmAddress) ZeroAsNative() AddressLike {
	if a.IsZero() {
		return EvmNative
	}
	return a
}

func (a *EvmAddress) ToBuffer() []byte {
	return a.inner.Bytes()
}

func (a *EvmAddress) ToHex() string {
	return a.inner.Hex()
}

func (a *EvmAddress) Equal(other AddressLike) bool {
	return strings.EqualFold(a.ToString(), other.ToString())
}

func (a *EvmAddress) IsNative() bool {
	return strings.EqualFold(a.inner.Hex(), NativeCurrency)
}

func (a *EvmAddress) IsZero() bool {
	return a.inner == common.HexToAddress(ZeroAddress)
}

func (a *EvmAddress) ToBigInt() *big.Int {
	return new(big.Int).SetBytes(a.inner.Bytes())
}

func (a *EvmAddress) ToString() string {
	return a.inner.Hex()
}

func (a *EvmAddress) ToJSON() string {
	return a.ToString()
}

func (a *EvmAddress) Inner() common.Address {
	return a.inner
}
