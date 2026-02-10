package auction

import (
	"encoding/hex"
	"math/big"

	"github.com/dawitel/cross-chain-sdk-go/crypto"
)

type AuctionDetails struct {
	StartTime       *big.Int
	InitialRateBump int
	Duration        *big.Int
	Points          []AuctionPoint
	GasCost         GasCost
}

type AuctionPoint struct {
	Delay         int
	ToTokenAmount string
}

type GasCost struct {
	GasBumpEstimate  *big.Int
	GasPriceEstimate *big.Int
}

func (ad *AuctionDetails) HashForSolana() ([]byte, error) {
	data, err := ad.encode()
	if err != nil {
		return nil, err
	}
	return crypto.Keccak256(data), nil
}

func (ad *AuctionDetails) encode() ([]byte, error) {
	var result []byte

	startTimeBytes := make([]byte, 8)
	ad.StartTime.FillBytes(startTimeBytes)
	result = append(result, startTimeBytes...)

	durationBytes := make([]byte, 8)
	ad.Duration.FillBytes(durationBytes)
	result = append(result, durationBytes...)

	initialRateBumpBytes := make([]byte, 2)
	initialRateBumpBytes[0] = byte(ad.InitialRateBump & 0xff)
	initialRateBumpBytes[1] = byte((ad.InitialRateBump >> 8) & 0xff)
	result = append(result, initialRateBumpBytes...)

	pointsCount := len(ad.Points)
	result = append(result, byte(pointsCount&0xff), byte((pointsCount>>8)&0xff))

	for _, point := range ad.Points {
		delayBytes := make([]byte, 4)
		big.NewInt(int64(point.Delay)).FillBytes(delayBytes)
		result = append(result, delayBytes...)

		amount, ok := new(big.Int).SetString(point.ToTokenAmount, 10)
		if !ok {
			return nil, &AuctionError{Message: "invalid toTokenAmount"}
		}
		amountBytes := make([]byte, 32)
		amount.FillBytes(amountBytes)
		result = append(result, amountBytes...)
	}

	gasBumpBytes := make([]byte, 8)
	ad.GasCost.GasBumpEstimate.FillBytes(gasBumpBytes)
	result = append(result, gasBumpBytes...)

	gasPriceBytes := make([]byte, 32)
	ad.GasCost.GasPriceEstimate.FillBytes(gasPriceBytes)
	result = append(result, gasPriceBytes...)

	return result, nil
}

func (ad *AuctionDetails) HashForSolanaHex() (string, error) {
	hash, err := ad.HashForSolana()
	if err != nil {
		return "", err
	}
	return "0x" + hex.EncodeToString(hash), nil
}

type AuctionError struct {
	Message string
}

func (e *AuctionError) Error() string {
	return e.Message
}
