package chains

type ChainType string

const (
	ChainTypeEVM ChainType = "EVM"
	ChainTypeSVM ChainType = "Solana"
)

func GetChainType(chain SupportedChain) ChainType {
	if IsSolana(chain) {
		return ChainTypeSVM
	}
	return ChainTypeEVM
}
