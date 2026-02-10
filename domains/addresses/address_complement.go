package addresses

import (
	"fmt"
	"math/big"
)

const UINT_160_MAX = "0xffffffffffffffffffffffffffffffffffffffff"

var UINT_160_MAX_BIG = mustBigIntFromHex(UINT_160_MAX)

type AddressComplement struct {
	inner *big.Int
}

var ZeroComplement = NewAddressComplement(big.NewInt(0))

func NewAddressComplement(val *big.Int) *AddressComplement {
	if val.Cmp(UINT_160_MAX_BIG) > 0 {
		panic(fmt.Sprintf("value exceeds UINT_160_MAX: %s", val.String()))
	}
	return &AddressComplement{inner: new(big.Int).Set(val)}
}

func (ac *AddressComplement) AsHex() string {
	hexStr := ac.inner.Text(16)
	if len(hexStr)%2 != 0 {
		hexStr = "0" + hexStr
	}
	return "0x" + hexStr
}

func (ac *AddressComplement) IsZero() bool {
	return ac.inner.Sign() == 0
}

func (ac *AddressComplement) Inner() *big.Int {
	return new(big.Int).Set(ac.inner)
}

func mustBigIntFromHex(hexStr string) *big.Int {
	hexStr = hexStr[2:]
	val, ok := new(big.Int).SetString(hexStr, 16)
	if !ok {
		panic("invalid hex")
	}
	return val
}
