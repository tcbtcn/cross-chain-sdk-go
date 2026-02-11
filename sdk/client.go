package sdk

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"math/big"
	"time"

	"github.com/dawitel/cross-chain-sdk-go/api/http"
	apiorders "github.com/dawitel/cross-chain-sdk-go/api/orders"
	"github.com/dawitel/cross-chain-sdk-go/api/quoter"
	"github.com/dawitel/cross-chain-sdk-go/api/relayer"
	"github.com/dawitel/cross-chain-sdk-go/chains"
	"github.com/dawitel/cross-chain-sdk-go/crypto"
	"github.com/dawitel/cross-chain-sdk-go/crypto/eip712"
	"github.com/dawitel/cross-chain-sdk-go/domains/addresses"
	"github.com/dawitel/cross-chain-sdk-go/domains/auction"
	"github.com/dawitel/cross-chain-sdk-go/domains/timelocks"
	"github.com/dawitel/cross-chain-sdk-go/orders"
)

type SDK struct {
	config     Config
	httpClient http.Client
	quoterApi  *quoter.QuoterApi
	ordersApi  *apiorders.OrdersApi
	relayerApi *relayer.RelayerApi
}

func NewSDK(config Config) *SDK {
	var httpClient http.Client = http.NewHTTPClient(config.URL, config.AuthKey)
	if config.HTTPClient != nil {
		if hc, ok := config.HTTPClient.(http.Client); ok {
			httpClient = hc
		}
	}

	return &SDK{
		config:     config,
		httpClient: httpClient,
		quoterApi: quoter.NewQuoterApi(quoter.QuoterApiConfig{
			URL:     config.URL,
			AuthKey: config.AuthKey,
		}, httpClient),
		ordersApi: apiorders.NewOrdersApi(apiorders.OrdersApiConfig{
			URL:     config.URL,
			AuthKey: config.AuthKey,
		}, httpClient),
		relayerApi: relayer.NewRelayerApi(relayer.RelayerApiConfig{
			URL:     config.URL,
			AuthKey: config.AuthKey,
		}, httpClient),
	}
}

func (s *SDK) GetQuote(ctx context.Context, params QuoteParams) (*quoter.QuoterResponse, error) {
	walletAddress := params.WalletAddress
	// Default walletAddress to zero address if empty for EVM requests (matching TypeScript SDK)
	if walletAddress == "" && chains.IsEvm(params.SrcChainID) {
		walletAddress = addresses.ZeroAddress
	}

	reqParams := quoter.QuoterRequestParams{
		SrcChain:        params.SrcChainID,
		DstChain:        params.DstChainID,
		SrcTokenAddress: params.SrcTokenAddress,
		DstTokenAddress: params.DstTokenAddress,
		Amount:          params.Amount,
		WalletAddress:   walletAddress,
		EnableEstimate:  params.EnableEstimate,
		Permit:          params.Permit,
		Fee:             params.TakingFeeBps,
		Source:          params.Source,
		IsPermit2:       params.IsPermit2,
	}

	return s.quoterApi.GetQuote(ctx, reqParams)
}

func (s *SDK) GetQuoteWithCustomPreset(ctx context.Context, params QuoteParams, customPreset quoter.CustomPreset) (*quoter.QuoterResponse, error) {
	walletAddress := params.WalletAddress
	// Default walletAddress to zero address if empty for EVM requests (matching TypeScript SDK)
	if walletAddress == "" && chains.IsEvm(params.SrcChainID) {
		walletAddress = addresses.ZeroAddress
	}

	reqParams := quoter.QuoterRequestParams{
		SrcChain:        params.SrcChainID,
		DstChain:        params.DstChainID,
		SrcTokenAddress: params.SrcTokenAddress,
		DstTokenAddress: params.DstTokenAddress,
		Amount:          params.Amount,
		WalletAddress:   walletAddress,
		EnableEstimate:  params.EnableEstimate,
		Permit:          params.Permit,
		Fee:             params.TakingFeeBps,
		Source:          params.Source,
		IsPermit2:       params.IsPermit2,
	}

	return s.quoterApi.GetQuoteWithCustomPreset(ctx, reqParams, customPreset)
}

func (s *SDK) CreateOrder(ctx context.Context, quote *quoter.QuoterResponse, params OrderParams) (*PreparedOrder, error) {
	if quote.QuoteID == "" {
		return nil, fmt.Errorf("request quote with enableEstimate=true")
	}

	order, hash, err := s.quoteToOrder(quote, params)
	if err != nil {
		return nil, err
	}

	return &PreparedOrder{
		Order:   order,
		Hash:    hash,
		QuoteID: quote.QuoteID,
	}, nil
}

func (s *SDK) SubmitOrder(ctx context.Context, srcChainID chains.SupportedChain, order *orders.EvmCrossChainOrder, quoteID string, secretHashes []string) (*OrderInfo, error) {
	signature, err := s.SignOrder(ctx, order, srcChainID)
	if err != nil {
		return nil, err
	}

	return s.submitOrder(ctx, srcChainID, order, quoteID, signature, secretHashes)
}

func (s *SDK) SubmitNativeOrder(ctx context.Context, srcChainID chains.SupportedChain, order *orders.EvmCrossChainOrder, maker *addresses.EvmAddress, quoteID string, secretHashes []string) (*OrderInfo, error) {
	signature := s.SignNativeOrder(order, maker)
	return s.submitOrder(ctx, srcChainID, order, quoteID, signature, secretHashes)
}

func (s *SDK) SignOrder(ctx context.Context, order *orders.EvmCrossChainOrder, srcChainID chains.SupportedChain) (string, error) {
	if s.config.BlockchainProvider == nil {
		return "", fmt.Errorf("blockchainProvider has not been set to config")
	}

	typedData, err := order.GetTypedData(srcChainID)
	if err != nil {
		return "", fmt.Errorf("failed to get typed data: %w", err)
	}
	typedDataMap, ok := typedData.(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid typed data format")
	}

	// Convert types from []map[string]string to []interface{}
	typesMap := typedDataMap["types"].(map[string]interface{})
	convertedTypes := make(map[string]interface{})
	for typeName, typeValue := range typesMap {
		if typeSlice, ok := typeValue.([]map[string]string); ok {
			convertedSlice := make([]interface{}, len(typeSlice))
			for i, v := range typeSlice {
				convertedSlice[i] = v
			}
			convertedTypes[typeName] = convertedSlice
		} else {
			convertedTypes[typeName] = typeValue
		}
	}

	eip712Data := eip712.TypedData{
		Types:       convertedTypes,
		PrimaryType: typedDataMap["primaryType"].(string),
		Domain:      typedDataMap["domain"].(map[string]interface{}),
		Message:     typedDataMap["message"].(map[string]interface{}),
	}

	orderStruct, err := order.Build()
	if err != nil {
		return "", fmt.Errorf("failed to build order: %w", err)
	}
	if orderStruct.Maker == "" {
		return "", fmt.Errorf("order maker address is empty")
	}
	return s.config.BlockchainProvider.SignTypedData(ctx, orderStruct.Maker, eip712Data)
}

func (s *SDK) SignNativeOrder(order *orders.EvmCrossChainOrder, maker *addresses.EvmAddress) string {
	// Use order's native signature method (matching TypeScript SDK)
	return order.NativeSignature(maker)
}

func (s *SDK) AnnounceOrder(ctx context.Context, order *orders.SolanaCrossChainOrder, quoteID string, secretHashes []string) (string, error) {
	if !order.MultipleFillsAllowed && len(secretHashes) > 1 {
		return "", fmt.Errorf("with disabled multiple fills you provided secretHashes > 1")
	} else if order.MultipleFillsAllowed && len(secretHashes) > 1 {
		partsCount, err := order.HashLock.GetPartsCount()
		if err != nil {
			return "", err
		}
		expectedCount := int(partsCount.Int64()) + 1
		if len(secretHashes) != expectedCount {
			return "", fmt.Errorf("secretHashes length should be equal to number of secrets: expected %d, got %d", expectedCount, len(secretHashes))
		}
	}

	// Use auction hash for Solana orders (matching TypeScript SDK)
	auctionHash, err := order.Auction.HashForSolanaHex()
	if err != nil {
		return "", fmt.Errorf("failed to get auction hash: %w", err)
	}

	hash := order.GetOrderHash()
	var secretHashesToSend []string
	if len(secretHashes) == 1 {
		secretHashesToSend = nil
	} else {
		secretHashesToSend = secretHashes
	}

	req := relayer.RelayerRequestSvm{
		Order:            order.ToJSON(),
		AuctionOrderHash: auctionHash,
		QuoteID:          quoteID,
		SecretHashes:     secretHashesToSend,
	}

	if err := s.relayerApi.Submit(ctx, req); err != nil {
		return "", err
	}

	return hash, nil
}

func (s *SDK) GetOrderStatus(ctx context.Context, orderHash string) (*apiorders.OrderStatusResponse, error) {
	return s.ordersApi.GetOrderStatus(ctx, orderHash)
}

func (s *SDK) GetActiveOrders(ctx context.Context, params apiorders.ActiveOrdersRequestParams) (*apiorders.ActiveOrdersResponse, error) {
	return s.ordersApi.GetActiveOrders(ctx, params)
}

func (s *SDK) GetOrdersByMaker(ctx context.Context, params apiorders.OrdersByMakerParams) (*apiorders.OrdersByMakerResponse, error) {
	return s.ordersApi.GetOrdersByMaker(ctx, params)
}

func (s *SDK) GetReadyToAcceptSecretFills(ctx context.Context, orderHash string) (*apiorders.ReadyToAcceptSecretFills, error) {
	return s.ordersApi.GetReadyToAcceptSecretFills(ctx, orderHash)
}

func (s *SDK) SubmitSecret(ctx context.Context, orderHash, secret string) error {
	return s.relayerApi.SubmitSecret(ctx, orderHash, secret)
}

func (s *SDK) GetPublishedSecrets(ctx context.Context, orderHash string) (*apiorders.PublishedSecretsResponse, error) {
	return s.ordersApi.GetPublishedSecrets(ctx, orderHash)
}

func (s *SDK) GetCancellableOrders(ctx context.Context, chainType chains.ChainType, page, limit int) (interface{}, error) {
	response, err := s.ordersApi.GetCancellableOrders(ctx, chainType, page, limit)
	if err != nil {
		return nil, err
	}

	// Transform response to typed structures (matching TypeScript SDK)
	if chainType == chains.ChainTypeEVM {
		return s.transformEvmCancellableOrders(response)
	}
	return s.transformSvmCancellableOrders(response)
}

func (s *SDK) transformEvmCancellableOrders(response *apiorders.CancellableOrdersResponse) (interface{}, error) {
	items := make([]EvmOrderCancellationData, 0, len(response.Items))

	for _, item := range response.Items {
		// Type assert to map to access fields
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			continue // Skip invalid items
		}

		// Check if it's an EVM order (has extension and srcChainId fields)
		if _, hasExtension := itemMap["extension"]; !hasExtension {
			continue // Not an EVM order
		}

		orderHash, _ := itemMap["orderHash"].(string)
		makerStr, _ := itemMap["maker"].(string)
		srcChainID, _ := itemMap["srcChainId"].(float64)
		dstChainID, _ := itemMap["dstChainId"].(float64)
		extension, _ := itemMap["extension"].(string)
		remainingMakerAmount, _ := itemMap["remainingMakerAmount"].(string)

		maker, err := addresses.EvmAddressFromString(makerStr)
		if err != nil {
			continue // Skip invalid maker address
		}

		// Parse order struct
		orderMap, ok := itemMap["order"].(map[string]interface{})
		if !ok {
			continue
		}

		order := orders.LimitOrderV4Struct{
			Maker:         getString(orderMap, "maker"),
			MakerAsset:    getString(orderMap, "makerAsset"),
			TakerAsset:    getString(orderMap, "takerAsset"),
			MakingAmount:  getString(orderMap, "makingAmount"),
			TakingAmount:  getString(orderMap, "takingAmount"),
			Receiver:      getString(orderMap, "receiver"),
			AllowedSender: getString(orderMap, "allowedSender"),
			MakerTraits:   getString(orderMap, "makerTraits"),
			Salt:          getString(orderMap, "salt"),
			Expiration:    getString(orderMap, "expiration"),
			Nonce:         getString(orderMap, "nonce"),
		}

		remainingAmount := big.NewInt(0)
		if remainingMakerAmount != "" {
			if amt, ok := new(big.Int).SetString(remainingMakerAmount, 10); ok {
				remainingAmount = amt
			}
		}

		items = append(items, EvmOrderCancellationData{
			OrderHash:            orderHash,
			Maker:                maker,
			SrcChainID:           chains.SupportedChain(srcChainID),
			DstChainID:           chains.SupportedChain(dstChainID),
			Order:                order,
			Extension:            extension,
			RemainingMakerAmount: remainingAmount,
		})
	}

	return struct {
		Items      []EvmOrderCancellationData `json:"items"`
		TotalCount int                        `json:"totalCount"`
		Page       int                        `json:"page"`
		Limit      int                        `json:"limit"`
	}{
		Items:      items,
		TotalCount: response.TotalCount,
		Page:       response.Page,
		Limit:      response.Limit,
	}, nil
}

func (s *SDK) transformSvmCancellableOrders(response *apiorders.CancellableOrdersResponse) (interface{}, error) {
	items := make([]SvmOrderCancellationData, 0, len(response.Items))

	for _, item := range response.Items {
		// Type assert to map to access fields
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			continue // Skip invalid items
		}

		// Check if it's an SVM order (has txSignature or no extension)
		if _, hasExtension := itemMap["extension"]; hasExtension {
			continue // Not an SVM order
		}

		orderHashStr, _ := itemMap["orderHash"].(string)
		makerStr, _ := itemMap["maker"].(string)

		maker, err := addresses.SolanaAddressFromString(makerStr)
		if err != nil {
			continue // Skip invalid maker address
		}

		// Parse order struct
		orderMap, ok := itemMap["order"].(map[string]interface{})
		if !ok {
			continue
		}

		orderInfoMap, ok := orderMap["orderInfo"].(map[string]interface{})
		if !ok {
			continue
		}

		tokenStr, _ := orderInfoMap["srcToken"].(string)
		token, err := addresses.SolanaAddressFromString(tokenStr)
		if err != nil {
			continue
		}

		extraMap, ok := orderMap["extra"].(map[string]interface{})
		if !ok {
			continue
		}

		srcAssetIsNative, _ := extraMap["srcAssetIsNative"].(bool)

		// Decode orderHash from base58 (Solana uses base58)
		// For now, we'll store as string and convert when needed
		orderHash := []byte(orderHashStr) // This should be base58 decoded, but for now use string bytes

		// Parse cancellation config
		cancellationConfigMap, _ := extraMap["resolverCancellationConfig"].(map[string]interface{})
		cancellationConfig := cancellationConfigMap // Store as interface for now

		items = append(items, SvmOrderCancellationData{
			OrderHash:          orderHash,
			Maker:              maker,
			Token:              token,
			CancellationConfig: cancellationConfig,
			IsAssetNative:      srcAssetIsNative,
		})
	}

	return struct {
		Items      []SvmOrderCancellationData `json:"items"`
		TotalCount int                        `json:"totalCount"`
		Page       int                        `json:"page"`
		Limit      int                        `json:"limit"`
	}{
		Items:      items,
		TotalCount: response.TotalCount,
		Page:       response.Page,
		Limit:      response.Limit,
	}, nil
}

// Helper function to safely get string from map
func getString(m map[string]interface{}, key string) string {
	if val, ok := m[key].(string); ok {
		return val
	}
	return ""
}

func (s *SDK) GetReadyToExecutePublicActions(ctx context.Context) (*apiorders.ReadyToExecutePublicActions, error) {
	return s.ordersApi.GetReadyToExecutePublicActions(ctx)
}

func (s *SDK) SubmitBatch(ctx context.Context, requests []interface{}) error {
	return s.relayerApi.SubmitBatch(ctx, requests)
}

func (s *SDK) BuildCancelOrderCallData(ctx context.Context, orderHash string) (string, error) {
	orderData, err := s.GetOrderStatus(ctx, orderHash)
	if err != nil {
		return "", fmt.Errorf("can not get order with the specified orderHash %s: %w", orderHash, err)
	}

	if !chains.IsEvm(orderData.SrcChainID) {
		return "", fmt.Errorf("expected evm src chain")
	}

	if orderData.Order == nil {
		return "", fmt.Errorf("order data not available in response")
	}

	makerTraits, ok := new(big.Int).SetString(orderData.Order.MakerTraits, 10)
	if !ok {
		makerTraits = big.NewInt(0)
	}

	return crypto.EncodeCancelOrder(orderHash, makerTraits)
}

func (s *SDK) PlaceOrder(ctx context.Context, quote *quoter.QuoterResponse, params OrderParams) (*OrderInfo, error) {
	preparedOrder, err := s.CreateOrder(ctx, quote, params)
	if err != nil {
		return nil, err
	}

	evmOrder, ok := preparedOrder.Order.(*orders.EvmCrossChainOrder)
	if !ok {
		return nil, fmt.Errorf("solana order must be announced with AnnounceOrder and placed onchain")
	}

	return s.SubmitOrder(ctx, chains.SupportedChain(quote.SrcChainID), evmOrder, preparedOrder.QuoteID, params.SecretHashes)
}

func (s *SDK) submitOrder(ctx context.Context, srcChainID chains.SupportedChain, order *orders.EvmCrossChainOrder, quoteID string, signature string, secretHashes []string) (*OrderInfo, error) {
	if !order.MultipleFillsAllowed && len(secretHashes) > 1 {
		return nil, fmt.Errorf("with disabled multiple fills you provided secretHashes > 1")
	} else if order.MultipleFillsAllowed && len(secretHashes) > 1 {
		partsCount, err := order.HashLock.GetPartsCount()
		if err != nil {
			return nil, err
		}
		expectedCount := int(partsCount.Int64()) + 1
		if len(secretHashes) != expectedCount {
			return nil, fmt.Errorf("secretHashes length should be equal to number of secrets: expected %d, got %d", expectedCount, len(secretHashes))
		}
	}

	orderStruct, err := order.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build order: %w", err)
	}

	var secretHashesToSend []string
	if len(secretHashes) == 1 {
		secretHashesToSend = nil
	} else {
		secretHashesToSend = secretHashes
	}

	req := relayer.RelayerRequestEvm{
		SrcChainID:   int(srcChainID),
		Order:        orderStruct,
		Signature:    signature,
		QuoteID:      quoteID,
		Extension:    order.Extension,
		SecretHashes: secretHashesToSend,
	}

	if err := s.relayerApi.Submit(ctx, req); err != nil {
		return nil, err
	}

	hash, err := order.GetOrderHash(srcChainID)
	if err != nil {
		return nil, fmt.Errorf("failed to get order hash: %w", err)
	}

	return &OrderInfo{
		Order:     orderStruct,
		Signature: signature,
		QuoteID:   quoteID,
		OrderHash: hash,
		Extension: order.Extension,
	}, nil
}

func (s *SDK) quoteToOrder(quote *quoter.QuoterResponse, params OrderParams) (interface{}, string, error) {
	if chains.IsEvm(chains.SupportedChain(quote.SrcChainID)) {
		return s.createEvmOrder(quote, params)
	}
	return s.createSolanaOrder(quote, params)
}

func (s *SDK) createEvmOrder(quote *quoter.QuoterResponse, params OrderParams) (*orders.EvmCrossChainOrder, string, error) {
	// Get preset - support custom preset (matching TypeScript SDK)
	preset := quote.Presets.Fast
	if params.CustomPreset != nil {
		// Custom preset is provided, use it
		if quote.Presets.Custom != nil {
			preset = *quote.Presets.Custom
		} else {
			return nil, "", fmt.Errorf("custom preset requested but not available in quote")
		}
	} else if params.Preset == "medium" {
		preset = quote.Presets.Medium
	} else if params.Preset == "slow" {
		preset = quote.Presets.Slow
	}

	makerAddr, err := addresses.EvmAddressFromString(params.WalletAddress)
	if err != nil {
		return nil, "", fmt.Errorf("invalid wallet address: %w", err)
	}
	if makerAddr == nil {
		return nil, "", fmt.Errorf("maker address is nil")
	}

	makerAsset, err := addresses.EvmAddressFromString(quote.SrcTokenAddress)
	if err != nil {
		return nil, "", fmt.Errorf("invalid source token address: %w", err)
	}
	if makerAsset == nil {
		return nil, "", fmt.Errorf("maker asset address is nil")
	}

	// Check if native asset (matching TypeScript SDK)
	isNativeAsset := makerAsset.IsNative()

	takerAsset, err := addresses.EvmAddressFromString(quote.DstTokenAddress)
	if err != nil {
		return nil, "", fmt.Errorf("invalid destination token address: %w", err)
	}
	if takerAsset == nil {
		return nil, "", fmt.Errorf("taker asset address is nil")
	}

	// Use zeroAsNative for takerAsset (matching TypeScript SDK)
	takerAsset = takerAsset.ZeroAsNative().(*addresses.EvmAddress)

	receiver := makerAddr
	if params.Receiver != "" {
		receiver, err = addresses.EvmAddressFromString(params.Receiver)
		if err != nil {
			return nil, "", fmt.Errorf("invalid receiver address: %w", err)
		}
		if receiver == nil {
			return nil, "", fmt.Errorf("receiver address is nil")
		}
	}

	makingAmount, ok := new(big.Int).SetString(quote.SrcTokenAmount, 10)
	if !ok {
		return nil, "", fmt.Errorf("invalid making amount")
	}

	takingAmount, ok := new(big.Int).SetString(preset.AuctionEndAmount, 10)
	if !ok {
		return nil, "", fmt.Errorf("invalid taking amount")
	}

	// Nonce generation logic matching TypeScript SDK
	// Nonce is required if !allowPartialFills || !allowMultipleFills
	allowPartialFills := preset.AllowPartialFills
	allowMultipleFills := preset.AllowMultipleFills
	isNonceRequired := !allowPartialFills || !allowMultipleFills

	var nonce *big.Int
	if isNonceRequired {
		if params.Nonce != nil {
			nonce = params.Nonce
		} else {
			// Generate random nonce (UINT_40_MAX = 2^40 - 1)
			nonceBytes := make([]byte, 5) // 40 bits = 5 bytes
			if _, err := rand.Read(nonceBytes); err != nil {
				return nil, "", fmt.Errorf("failed to generate nonce: %w", err)
			}
			nonce = new(big.Int).SetBytes(nonceBytes)
			// Ensure it fits in 40 bits
			maxNonce := new(big.Int).Lsh(big.NewInt(1), 40)
			maxNonce.Sub(maxNonce, big.NewInt(1))
			nonce.Mod(nonce, maxNonce)
		}
	} else {
		nonce = params.Nonce // Can be nil if not required
	}

	salt := big.NewInt(time.Now().UnixNano())
	deadline := big.NewInt(time.Now().Unix() + int64(preset.AuctionDuration))

	tl, err := timelocks.NewTimeLocks(timelocks.TimeLocksParams{
		SrcWithdrawal:         big.NewInt(int64(quote.TimeLocks.SrcWithdrawal)),
		SrcPublicWithdrawal:   big.NewInt(int64(quote.TimeLocks.SrcPublicWithdrawal)),
		SrcCancellation:       big.NewInt(int64(quote.TimeLocks.SrcCancellation)),
		SrcPublicCancellation: big.NewInt(int64(quote.TimeLocks.SrcPublicCancellation)),
		DstWithdrawal:         big.NewInt(int64(quote.TimeLocks.DstWithdrawal)),
		DstPublicWithdrawal:   big.NewInt(int64(quote.TimeLocks.DstPublicWithdrawal)),
		DstCancellation:       big.NewInt(int64(quote.TimeLocks.DstCancellation)),
	})
	if err != nil {
		return nil, "", err
	}

	srcSafetyDeposit, _ := new(big.Int).SetString(quote.SrcSafetyDeposit, 10)
	dstSafetyDeposit, _ := new(big.Int).SetString(quote.DstSafetyDeposit, 10)

	if params.HashLock == nil {
		return nil, "", fmt.Errorf("hashLock is required")
	}

	// Build whitelist (matching TypeScript SDK getWhitelist method)
	whitelist := s.buildWhitelist(quote.Whitelist, big.NewInt(time.Now().Unix()), preset.ExclusiveResolver)

	// Create auction details
	gasPriceEstimate := big.NewInt(0)
	if preset.GasCost.GasPriceEstimate != "" {
		if gp, ok := new(big.Int).SetString(preset.GasCost.GasPriceEstimate, 10); ok {
			gasPriceEstimate = gp
		}
	}

	auctionDetails := &auction.AuctionDetails{
		StartTime:       big.NewInt(time.Now().Unix()),
		Duration:        big.NewInt(int64(preset.AuctionDuration)),
		InitialRateBump: preset.InitialRateBump,
		Points:          preset.Points,
		GasCost: auction.GasCost{
			GasBumpEstimate:  big.NewInt(int64(preset.GasCost.GasBumpEstimate)),
			GasPriceEstimate: gasPriceEstimate,
		},
	}

	// Create SettlementPostInteractionData
	postInteractionData := orders.NewSettlementPostInteractionData(
		whitelist,
		big.NewInt(time.Now().Unix()),
	)

	// Determine AddressComplement (for EVM addresses, complement is always zero)
	dstAddressFirstPart := addresses.ZeroComplement

	// Get escrow factory address
	escrowFactory, err := addresses.EvmAddressFromString(quote.SrcEscrowFactory)
	if err != nil {
		return nil, "", fmt.Errorf("invalid escrow factory address: %w", err)
	}

	// Create EscrowExtension
	escrowExt := orders.NewEscrowExtension(
		escrowFactory,
		auctionDetails,
		postInteractionData,
		nil, // makerPermit (optional, not implemented yet)
		params.HashLock,
		chains.SupportedChain(quote.DstChainID),
		takerAsset,
		srcSafetyDeposit,
		dstSafetyDeposit,
		tl,
		dstAddressFirstPart,
	)

	// Build extension
	extension, err := escrowExt.Build()
	if err != nil {
		return nil, "", fmt.Errorf("failed to build extension: %w", err)
	}

	// For native orders, we need to use fromNative equivalent
	// For now, we'll create a regular order but mark it appropriately
	// TODO: Implement proper native order creation when NativeOrderFactory is available
	if isNativeAsset {
		if quote.NativeOrderFactory == "" || quote.NativeOrderImpl == "" {
			return nil, "", fmt.Errorf("native order factory not available in quote for native asset order")
		}
		// Native orders require special handling - for now we'll proceed with regular order
		// but this should be properly implemented with fromNative equivalent
	}

	order := &orders.EvmCrossChainOrder{
		Maker:                makerAddr,
		MakerAsset:           makerAsset,
		TakerAsset:           takerAsset,
		MakingAmount:         makingAmount,
		TakingAmount:         takingAmount,
		Receiver:             receiver,
		Deadline:             deadline,
		AuctionStartTime:     big.NewInt(time.Now().Unix()),
		AuctionEndTime:       big.NewInt(time.Now().Unix() + int64(preset.AuctionDuration)),
		Nonce:                nonce,
		Salt:                 salt,
		HashLock:             params.HashLock,
		TimeLocks:            tl,
		SrcSafetyDeposit:     srcSafetyDeposit,
		DstSafetyDeposit:     dstSafetyDeposit,
		DstChainID:           chains.SupportedChain(quote.DstChainID),
		Extension:            extension,
		MultipleFillsAllowed: preset.AllowMultipleFills,
		Whitelist:            whitelist,
	}

	hash, err := order.GetOrderHash(chains.SupportedChain(quote.SrcChainID))
	if err != nil {
		return nil, "", fmt.Errorf("failed to get order hash: %w", err)
	}

	return order, hash, nil
}

func (s *SDK) createSolanaOrder(quote *quoter.QuoterResponse, params OrderParams) (*orders.SolanaCrossChainOrder, string, error) {
	// Require receiver for Solana orders (matching TypeScript SDK)
	if params.Receiver == "" {
		return nil, "", fmt.Errorf("receiver is required for solana order")
	}

	receiver, err := addresses.EvmAddressFromString(params.Receiver)
	if err != nil {
		return nil, "", fmt.Errorf("invalid receiver address: %w", err)
	}
	if receiver == nil {
		return nil, "", fmt.Errorf("receiver address is nil")
	}

	// Get preset - support custom preset (matching TypeScript SDK)
	preset := quote.Presets.Fast
	if params.CustomPreset != nil {
		// Custom preset is provided, use it
		if quote.Presets.Custom != nil {
			preset = *quote.Presets.Custom
		} else {
			return nil, "", fmt.Errorf("custom preset requested but not available in quote")
		}
	} else if params.Preset == "medium" {
		preset = quote.Presets.Medium
	} else if params.Preset == "slow" {
		preset = quote.Presets.Slow
	}

	// Get source token address (Solana)
	srcTokenAddr, err := addresses.SolanaAddressFromString(quote.SrcTokenAddress)
	if err != nil {
		return nil, "", fmt.Errorf("invalid source token address: %w", err)
	}

	// Determine if source asset is native (matching TypeScript SDK)
	srcAssetIsNative := srcTokenAddr.IsNative()

	// Get destination token address (EVM)
	dstTokenAddr, err := addresses.EvmAddressFromString(quote.DstTokenAddress)
	if err != nil {
		return nil, "", fmt.Errorf("invalid destination token address: %w", err)
	}

	// Get maker address (Solana) - use wallet address from params
	if params.WalletAddress == "" {
		return nil, "", fmt.Errorf("wallet address is required for solana order")
	}
	makerAddr, err := addresses.SolanaAddressFromString(params.WalletAddress)
	if err != nil {
		return nil, "", fmt.Errorf("invalid maker address: %w", err)
	}

	// Parse amounts
	srcAmount, ok := new(big.Int).SetString(quote.SrcTokenAmount, 10)
	if !ok {
		return nil, "", fmt.Errorf("invalid source token amount")
	}

	takingAmount, ok := new(big.Int).SetString(preset.AuctionEndAmount, 10)
	if !ok {
		return nil, "", fmt.Errorf("invalid taking amount")
	}

	// Create time locks
	tl, err := timelocks.NewTimeLocks(timelocks.TimeLocksParams{
		SrcWithdrawal:         big.NewInt(int64(quote.TimeLocks.SrcWithdrawal)),
		SrcPublicWithdrawal:   big.NewInt(int64(quote.TimeLocks.SrcPublicWithdrawal)),
		SrcCancellation:       big.NewInt(int64(quote.TimeLocks.SrcCancellation)),
		SrcPublicCancellation: big.NewInt(int64(quote.TimeLocks.SrcPublicCancellation)),
		DstWithdrawal:         big.NewInt(int64(quote.TimeLocks.DstWithdrawal)),
		DstPublicWithdrawal:   big.NewInt(int64(quote.TimeLocks.DstPublicWithdrawal)),
		DstCancellation:       big.NewInt(int64(quote.TimeLocks.DstCancellation)),
	})
	if err != nil {
		return nil, "", err
	}

	// Parse safety deposits
	srcSafetyDeposit, _ := new(big.Int).SetString(quote.SrcSafetyDeposit, 10)
	dstSafetyDeposit, _ := new(big.Int).SetString(quote.DstSafetyDeposit, 10)

	// Parse gas price estimate
	gasPriceEstimate := big.NewInt(0)
	if preset.GasCost.GasPriceEstimate != "" {
		if gp, ok := new(big.Int).SetString(preset.GasCost.GasPriceEstimate, 10); ok {
			gasPriceEstimate = gp
		}
	}

	// Create auction details
	auctionDetails := &auction.AuctionDetails{
		StartTime:       big.NewInt(time.Now().Unix()),
		Duration:        big.NewInt(int64(preset.AuctionDuration)),
		InitialRateBump: preset.InitialRateBump,
		Points:          preset.Points,
		GasCost: auction.GasCost{
			GasBumpEstimate:  big.NewInt(int64(preset.GasCost.GasBumpEstimate)),
			GasPriceEstimate: gasPriceEstimate,
		},
	}

	// Generate order hash (will be computed from order data)
	orderHash := make([]byte, 32)
	if _, err := rand.Read(orderHash); err != nil {
		return nil, "", fmt.Errorf("failed to generate random order hash: %w", err)
	}

	if params.HashLock == nil {
		return nil, "", fmt.Errorf("hashLock is required")
	}

	order := &orders.SolanaCrossChainOrder{
		OrderHash:            orderHash,
		HashLock:             params.HashLock,
		Auction:              auctionDetails,
		MultipleFillsAllowed: preset.AllowMultipleFills,
		// Store additional data needed for ToJSON
		SrcToken:         srcTokenAddr,
		DstToken:         dstTokenAddr,
		Maker:            makerAddr,
		Receiver:         receiver,
		SrcAmount:        srcAmount,
		MinDstAmount:     takingAmount,
		SrcSafetyDeposit: srcSafetyDeposit,
		DstSafetyDeposit: dstSafetyDeposit,
		TimeLocks:        tl,
		DstChainID:       chains.SupportedChain(quote.DstChainID),
		Salt:             big.NewInt(time.Now().UnixNano()),
		Source:           params.Source,
		SrcAssetIsNative: srcAssetIsNative,
	}

	hash := order.GetOrderHash()
	return order, hash, nil
}

// buildWhitelist creates whitelist items from quote whitelist addresses
// Matching TypeScript SDK getWhitelist method
func (s *SDK) buildWhitelist(whitelistAddrs []string, auctionStartTime *big.Int, exclusiveResolver string) []orders.AuctionWhitelistItem {
	if len(whitelistAddrs) == 0 {
		return nil
	}

	items := make([]orders.AuctionWhitelistItem, 0, len(whitelistAddrs))
	for _, addrStr := range whitelistAddrs {
		addr, err := addresses.EvmAddressFromString(addrStr)
		if err != nil {
			continue // Skip invalid addresses
		}

		var allowFrom *big.Int
		if exclusiveResolver != "" {
			exclusiveAddr, err := addresses.EvmAddressFromString(exclusiveResolver)
			if err == nil && addr.Equal(exclusiveAddr) {
				allowFrom = big.NewInt(0) // Exclusive resolver can execute from start
			} else {
				// allowFrom is uint16, so we use 0 for non-exclusive resolvers
				// The actual timestamp check happens elsewhere
				allowFrom = big.NewInt(0)
			}
		} else {
			allowFrom = big.NewInt(0) // No exclusive resolver, all can execute from start
		}

		items = append(items, orders.AuctionWhitelistItem{
			Address:   addr,
			AllowFrom: allowFrom,
		})
	}

	return items
}

func GenerateSecrets(count int) ([]string, error) {
	secrets := make([]string, count)
	for i := 0; i < count; i++ {
		secret := make([]byte, 32)
		if _, err := rand.Read(secret); err != nil {
			return nil, err
		}
		secrets[i] = "0x" + hex.EncodeToString(secret)
	}
	return secrets, nil
}
