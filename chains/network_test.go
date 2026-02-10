package chains

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsSupportedChain(t *testing.T) {
	tests := []struct {
		name  string
		chain NetworkEnum
		want  bool
	}{
		{
			name:  "Ethereum is supported",
			chain: Ethereum,
			want:  true,
		},
		{
			name:  "Polygon is supported",
			chain: Polygon,
			want:  true,
		},
		{
			name:  "Binance is supported",
			chain: Binance,
			want:  true,
		},
		{
			name:  "Solana is supported",
			chain: Solana,
			want:  true,
		},
		{
			name:  "Arbitrum is supported",
			chain: Arbitrum,
			want:  true,
		},
		{
			name:  "Avalanche is supported",
			chain: Avalanche,
			want:  true,
		},
		{
			name:  "Optimism is supported",
			chain: Optimism,
			want:  true,
		},
		{
			name:  "Fantom is supported",
			chain: Fantom,
			want:  true,
		},
		{
			name:  "Gnosis is supported",
			chain: Gnosis,
			want:  true,
		},
		{
			name:  "Coinbase is supported",
			chain: Coinbase,
			want:  true,
		},
		{
			name:  "ZkSync is supported",
			chain: ZkSync,
			want:  true,
		},
		{
			name:  "Linea is supported",
			chain: Linea,
			want:  true,
		},
		{
			name:  "Sonic is supported",
			chain: Sonic,
			want:  true,
		},
		{
			name:  "Unichain is supported",
			chain: Unichain,
			want:  true,
		},
		{
			name:  "zero is not supported",
			chain: NetworkEnum(0),
			want:  false,
		},
		{
			name:  "random chain ID is not supported",
			chain: NetworkEnum(999999),
			want:  false,
		},
		{
			name:  "max uint64 is not supported",
			chain: NetworkEnum(^uint64(0)),
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsSupportedChain(tt.chain)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestIsEvm(t *testing.T) {
	tests := []struct {
		name  string
		chain SupportedChain
		want  bool
	}{
		{
			name:  "Ethereum is EVM",
			chain: Ethereum,
			want:  true,
		},
		{
			name:  "Polygon is EVM",
			chain: Polygon,
			want:  true,
		},
		{
			name:  "Binance is EVM",
			chain: Binance,
			want:  true,
		},
		{
			name:  "Arbitrum is EVM",
			chain: Arbitrum,
			want:  true,
		},
		{
			name:  "Solana is not EVM",
			chain: Solana,
			want:  false,
		},
		{
			name:  "unsupported chain is not EVM",
			chain: NetworkEnum(999999),
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsEvm(tt.chain)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestIsSolana(t *testing.T) {
	tests := []struct {
		name  string
		chain SupportedChain
		want  bool
	}{
		{
			name:  "Solana is Solana",
			chain: Solana,
			want:  true,
		},
		{
			name:  "Ethereum is not Solana",
			chain: Ethereum,
			want:  false,
		},
		{
			name:  "Polygon is not Solana",
			chain: Polygon,
			want:  false,
		},
		{
			name:  "unsupported chain is not Solana",
			chain: NetworkEnum(999999),
			want:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsSolana(tt.chain)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestSupportedChains_AllChains(t *testing.T) {
	expectedChains := map[NetworkEnum]bool{
		Ethereum:  true,
		Polygon:   true,
		Binance:   true,
		Optimism:  true,
		Arbitrum:  true,
		Avalanche: true,
		Fantom:    true,
		Gnosis:    true,
		Coinbase:  true,
		ZkSync:    true,
		Linea:     true,
		Sonic:     true,
		Unichain:  true,
		Solana:    true,
	}

	for _, chain := range SupportedChains {
		assert.True(t, expectedChains[chain], "chain %d should be in expected list", chain)
		assert.True(t, IsSupportedChain(chain), "chain %d should be supported", chain)
	}

	assert.Equal(t, 14, len(SupportedChains), "should have 14 supported chains")
}
