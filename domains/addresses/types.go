package addresses

import "math/big"

type AddressLike interface {
	ToString() string
	NativeAsZero() AddressLike
	ZeroAsNative() AddressLike
	ToBuffer() []byte
	ToHex() string
	Equal(other AddressLike) bool
	IsNative() bool
	IsZero() bool
	ToBigInt() *big.Int
}
