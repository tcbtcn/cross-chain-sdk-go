package chains

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetChainType(t *testing.T) {
	tests := []struct {
		name  string
		chain SupportedChain
		want  ChainType
	}{
		{
			name:  "Ethereum is EVM",
			chain: Ethereum,
			want:  ChainTypeEVM,
		},
		{
			name:  "Polygon is EVM",
			chain: Polygon,
			want:  ChainTypeEVM,
		},
		{
			name:  "Binance is EVM",
			chain: Binance,
			want:  ChainTypeEVM,
		},
		{
			name:  "Arbitrum is EVM",
			chain: Arbitrum,
			want:  ChainTypeEVM,
		},
		{
			name:  "Avalanche is EVM",
			chain: Avalanche,
			want:  ChainTypeEVM,
		},
		{
			name:  "Optimism is EVM",
			chain: Optimism,
			want:  ChainTypeEVM,
		},
		{
			name:  "Fantom is EVM",
			chain: Fantom,
			want:  ChainTypeEVM,
		},
		{
			name:  "Gnosis is EVM",
			chain: Gnosis,
			want:  ChainTypeEVM,
		},
		{
			name:  "Coinbase is EVM",
			chain: Coinbase,
			want:  ChainTypeEVM,
		},
		{
			name:  "ZkSync is EVM",
			chain: ZkSync,
			want:  ChainTypeEVM,
		},
		{
			name:  "Linea is EVM",
			chain: Linea,
			want:  ChainTypeEVM,
		},
		{
			name:  "Sonic is EVM",
			chain: Sonic,
			want:  ChainTypeEVM,
		},
		{
			name:  "Unichain is EVM",
			chain: Unichain,
			want:  ChainTypeEVM,
		},
		{
			name:  "Solana is SVM",
			chain: Solana,
			want:  ChainTypeSVM,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetChainType(tt.chain)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestChainTypeConstants(t *testing.T) {
	assert.Equal(t, "EVM", string(ChainTypeEVM))
	assert.Equal(t, "Solana", string(ChainTypeSVM))
}
