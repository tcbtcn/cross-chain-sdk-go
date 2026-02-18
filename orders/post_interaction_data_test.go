package orders

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tcbtcn/cross-chain-sdk-go/domains/addresses"
)

func TestSettlementPostInteractionData_Encode(t *testing.T) {
	whitelist := []AuctionWhitelistItem{
		{
			Address:   addresses.NewEvmAddress("0x1111111111111111111111111111111111111111"),
			AllowFrom: big.NewInt(1000),
		},
		{
			Address:   addresses.NewEvmAddress("0x2222222222222222222222222222222222222222"),
			AllowFrom: big.NewInt(2000),
		},
	}

	spid := NewSettlementPostInteractionData(whitelist, big.NewInt(1735689600))
	encoded, err := spid.Encode()
	require.NoError(t, err)
	assert.NotEmpty(t, encoded)
	assert.Equal(t, "0x", encoded[:2])
}

func TestSettlementPostInteractionData_Encode_EmptyWhitelist(t *testing.T) {
	spid := NewSettlementPostInteractionData([]AuctionWhitelistItem{}, big.NewInt(1735689600))
	encoded, err := spid.Encode()
	require.NoError(t, err)
	assert.NotEmpty(t, encoded)
}

func TestSettlementPostInteractionData_EncodeForExtension(t *testing.T) {
	whitelist := []AuctionWhitelistItem{
		{
			Address:   addresses.NewEvmAddress("0x1111111111111111111111111111111111111111"),
			AllowFrom: big.NewInt(1000),
		},
	}

	spid := NewSettlementPostInteractionData(whitelist, big.NewInt(1735689600))
	settlementAddr := addresses.NewEvmAddress("0x3333333333333333333333333333333333333333")

	encoded, err := spid.EncodeForExtension(settlementAddr)
	require.NoError(t, err)
	assert.NotEmpty(t, encoded)
	assert.Contains(t, encoded, settlementAddr.ToHex()[2:])
}

func TestSettlementPostInteractionData_EncodeForExtension_NilAddress(t *testing.T) {
	spid := NewSettlementPostInteractionData([]AuctionWhitelistItem{}, big.NewInt(1735689600))
	_, err := spid.EncodeForExtension(nil)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "settlement address is required")
}

func TestSettlementPostInteractionData_Encode_NilWhitelistItem(t *testing.T) {
	whitelist := []AuctionWhitelistItem{
		{
			Address:   nil,
			AllowFrom: big.NewInt(1000),
		},
	}

	spid := NewSettlementPostInteractionData(whitelist, big.NewInt(1735689600))
	_, err := spid.Encode()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "whitelist item address is nil")
}

func TestSettlementPostInteractionData_Encode_NilAllowFrom(t *testing.T) {
	whitelist := []AuctionWhitelistItem{
		{
			Address:   addresses.NewEvmAddress("0x1111111111111111111111111111111111111111"),
			AllowFrom: nil,
		},
	}

	spid := NewSettlementPostInteractionData(whitelist, big.NewInt(1735689600))
	_, err := spid.Encode()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "whitelist item allowFrom is nil")
}

func TestSettlementPostInteractionData_Encode_AllowFromExceedsUint16(t *testing.T) {
	whitelist := []AuctionWhitelistItem{
		{
			Address:   addresses.NewEvmAddress("0x1111111111111111111111111111111111111111"),
			AllowFrom: big.NewInt(0x10000),
		},
	}

	spid := NewSettlementPostInteractionData(whitelist, big.NewInt(1735689600))
	_, err := spid.Encode()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "allowFrom exceeds uint16 max")
}

func TestSettlementPostInteractionData_Encode_MultipleWhitelistItems(t *testing.T) {
	whitelist := []AuctionWhitelistItem{
		{
			Address:   addresses.NewEvmAddress("0x1111111111111111111111111111111111111111"),
			AllowFrom: big.NewInt(1000),
		},
		{
			Address:   addresses.NewEvmAddress("0x2222222222222222222222222222222222222222"),
			AllowFrom: big.NewInt(2000),
		},
		{
			Address:   addresses.NewEvmAddress("0x3333333333333333333333333333333333333333"),
			AllowFrom: big.NewInt(3000),
		},
	}

	spid := NewSettlementPostInteractionData(whitelist, big.NewInt(1735689600))
	encoded, err := spid.Encode()
	require.NoError(t, err)
	assert.NotEmpty(t, encoded)
}

func TestSettlementPostInteractionData_Encode_Deterministic(t *testing.T) {
	whitelist := []AuctionWhitelistItem{
		{
			Address:   addresses.NewEvmAddress("0x1111111111111111111111111111111111111111"),
			AllowFrom: big.NewInt(1000),
		},
	}

	spid := NewSettlementPostInteractionData(whitelist, big.NewInt(1735689600))
	encoded1, err := spid.Encode()
	require.NoError(t, err)

	encoded2, err := spid.Encode()
	require.NoError(t, err)

	assert.Equal(t, encoded1, encoded2)
}

func TestDecodeSettlementPostInteractionData(t *testing.T) {
	whitelist := []AuctionWhitelistItem{
		{
			Address:   addresses.NewEvmAddress("0x1111111111111111111111111111111111111111"),
			AllowFrom: big.NewInt(1000),
		},
	}

	spid := NewSettlementPostInteractionData(whitelist, big.NewInt(1735689600))
	encoded, err := spid.Encode()
	require.NoError(t, err)

	decoded, err := DecodeSettlementPostInteractionData(encoded)
	require.NoError(t, err)
	assert.NotNil(t, decoded)
	assert.Equal(t, spid.ResolvingStartTime.String(), decoded.ResolvingStartTime.String())
	assert.Equal(t, len(spid.Whitelist), len(decoded.Whitelist))
}
