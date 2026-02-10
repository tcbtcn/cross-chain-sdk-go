package addresses

import (
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/gagliardetto/solana-go"
	"github.com/mr-tron/base58"
)

var (
	SolanaAssociatedTokenProgramID = MustSolanaAddressFromString("ATokenGPvbdGVxr1b2hvZbsiqW5xWH25efTNsLJA8knL")
	SolanaTokenProgramID           = MustSolanaAddressFromString("TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA")
	SolanaToken2022ProgramID       = MustSolanaAddressFromString("TokenzQdBNbLqP5VEhdkAS6EPFLC1PHnBqCXEpPxuEb")
	SolanaSystemProgramID          = MustSolanaAddressFromString("11111111111111111111111111111111")
	SolanaSysvarRentID             = MustSolanaAddressFromString("SysvarRent111111111111111111111111111111111")
	SolanaWrappedNative            = MustSolanaAddressFromString("So11111111111111111111111111111111111111112")
	SolanaNative                   = MustSolanaAddressFromString("SoNative11111111111111111111111111111111111")
	SolanaZero                     = SolanaAddressFromBigInt(big.NewInt(0))
)

type SolanaAddress struct {
	buf []byte
}

func NewSolanaAddress(value string) (*SolanaAddress, error) {
	buf, err := base58.Decode(value)
	if err != nil {
		return nil, fmt.Errorf("%s is not a valid address: %w", value, err)
	}
	if len(buf) != 32 {
		return nil, fmt.Errorf("address must be 32 bytes, got %d", len(buf))
	}
	return &SolanaAddress{buf: buf}, nil
}

func MustSolanaAddressFromString(str string) *SolanaAddress {
	addr, err := NewSolanaAddress(str)
	if err != nil {
		panic(err)
	}
	return addr
}

func SolanaAddressFromString(str string) (*SolanaAddress, error) {
	return NewSolanaAddress(str)
}

func SolanaAddressFromBuffer(buf []byte) (*SolanaAddress, error) {
	if len(buf) != 32 {
		return nil, fmt.Errorf("buffer must be 32 bytes, got %d", len(buf))
	}
	addr := &SolanaAddress{buf: make([]byte, 32)}
	copy(addr.buf, buf)
	return addr, nil
}

func SolanaAddressFromPublicKey(pubkey solana.PublicKey) *SolanaAddress {
	return &SolanaAddress{buf: pubkey.Bytes()}
}

func SolanaAddressFromBigInt(val *big.Int) *SolanaAddress {
	hexStr := fmt.Sprintf("%064x", val)
	buf, _ := hex.DecodeString(hexStr)
	addr, _ := SolanaAddressFromBuffer(buf)
	return addr
}

func (a *SolanaAddress) NativeAsZero() AddressLike {
	return a
}

func (a *SolanaAddress) ZeroAsNative() AddressLike {
	return a
}

func (a *SolanaAddress) ToString() string {
	return base58.Encode(a.buf)
}

func (a *SolanaAddress) ToJSON() string {
	return a.ToString()
}

func (a *SolanaAddress) ToBuffer() []byte {
	buf := make([]byte, 32)
	copy(buf, a.buf)
	return buf
}

func (a *SolanaAddress) Equal(other AddressLike) bool {
	otherBuf := other.ToBuffer()
	if len(a.buf) != len(otherBuf) {
		return false
	}
	for i := range a.buf {
		if a.buf[i] != otherBuf[i] {
			return false
		}
	}
	return true
}

func (a *SolanaAddress) IsNative() bool {
	return a.Equal(SolanaNative)
}

func (a *SolanaAddress) IsZero() bool {
	return a.Equal(SolanaZero)
}

func (a *SolanaAddress) ToHex() string {
	return "0x" + hex.EncodeToString(a.buf)
}

func (a *SolanaAddress) ToBigInt() *big.Int {
	return new(big.Int).SetBytes(a.buf)
}

func (a *SolanaAddress) ToPublicKey() solana.PublicKey {
	var pubkey solana.PublicKey
	copy(pubkey[:], a.buf)
	return pubkey
}
