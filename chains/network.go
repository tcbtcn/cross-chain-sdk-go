package chains

type NetworkEnum uint64

const (
	Ethereum  NetworkEnum = 1
	Polygon   NetworkEnum = 137
	ZkSync    NetworkEnum = 324
	Binance   NetworkEnum = 56
	Arbitrum  NetworkEnum = 42161
	Avalanche NetworkEnum = 43114
	Optimism  NetworkEnum = 10
	Fantom    NetworkEnum = 250
	Gnosis    NetworkEnum = 100
	Coinbase  NetworkEnum = 8453
	Linea     NetworkEnum = 59144
	Sonic     NetworkEnum = 146
	Unichain  NetworkEnum = 130
	Solana    NetworkEnum = 501
)

var SupportedChains = []NetworkEnum{
	Ethereum,
	Polygon,
	Binance,
	Optimism,
	Arbitrum,
	Avalanche,
	Fantom,
	Gnosis,
	Coinbase,
	ZkSync,
	Linea,
	Sonic,
	Unichain,
	Solana,
}

type SupportedChain = NetworkEnum
type EvmChain = NetworkEnum
type SolanaChain = NetworkEnum

func IsSupportedChain(chain NetworkEnum) bool {
	for _, supported := range SupportedChains {
		if chain == supported {
			return true
		}
	}
	return false
}

func IsEvm(chain SupportedChain) bool {
	return IsSupportedChain(chain) && chain != Solana
}

func IsSolana(chain SupportedChain) bool {
	return chain == Solana
}
