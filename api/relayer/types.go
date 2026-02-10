package relayer

type RelayerApiConfig struct {
	URL     string
	AuthKey string
}

type RelayerRequestEvm struct {
	SrcChainID   int         `json:"srcChainId"`
	Order        interface{} `json:"order"`
	Signature    string      `json:"signature"`
	QuoteID      string      `json:"quoteId"`
	Extension    string      `json:"extension"`
	SecretHashes []string    `json:"secretHashes,omitempty"`
}

type RelayerRequestSvm struct {
	Order            interface{} `json:"order"`
	AuctionOrderHash string      `json:"auctionOrderHash"`
	QuoteID          string      `json:"quoteId"`
	SecretHashes     []string    `json:"secretHashes,omitempty"`
}

type SubmitSecretRequest struct {
	OrderHash string `json:"orderHash"`
	Secret    string `json:"secret"`
}
