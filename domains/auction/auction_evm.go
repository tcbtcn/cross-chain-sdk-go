package auction

import (
	"encoding/binary"
	"math/big"
)

func (ad *AuctionDetails) EncodeForEvm() ([]byte, error) {
	var result []byte

	gasBumpEstimate := ad.GasCost.GasBumpEstimate.Uint64()
	if gasBumpEstimate > 0xffffff {
		return nil, &AuctionError{Message: "gasBumpEstimate exceeds uint24"}
	}
	gasBumpBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(gasBumpBytes, uint32(gasBumpEstimate))
	result = append(result, gasBumpBytes[1:]...)

	gasPriceEstimate := ad.GasCost.GasPriceEstimate.Uint64()
	if gasPriceEstimate > 0xffffffff {
		return nil, &AuctionError{Message: "gasPriceEstimate exceeds uint32"}
	}
	gasPriceBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(gasPriceBytes, uint32(gasPriceEstimate))
	result = append(result, gasPriceBytes...)

	startTime := ad.StartTime.Uint64()
	if startTime > 0xffffffff {
		return nil, &AuctionError{Message: "startTime exceeds uint32"}
	}
	startTimeBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(startTimeBytes, uint32(startTime))
	result = append(result, startTimeBytes...)

	duration := ad.Duration.Uint64()
	if duration > 0xffffff {
		return nil, &AuctionError{Message: "duration exceeds uint24"}
	}
	durationBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(durationBytes, uint32(duration))
	result = append(result, durationBytes[1:]...)

	initialRateBump := uint64(ad.InitialRateBump)
	if initialRateBump > 0xffffff {
		return nil, &AuctionError{Message: "initialRateBump exceeds uint24"}
	}
	initialRateBumpBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(initialRateBumpBytes, uint32(initialRateBump))
	result = append(result, initialRateBumpBytes[1:]...)

	for _, point := range ad.Points {
		delay := uint64(point.Delay)
		if delay > 0xffffff {
			return nil, &AuctionError{Message: "point delay exceeds uint24"}
		}
		delayBytes := make([]byte, 4)
		binary.BigEndian.PutUint32(delayBytes, uint32(delay))
		result = append(result, delayBytes[1:]...)

		amount, ok := new(big.Int).SetString(point.ToTokenAmount, 10)
		if !ok {
			return nil, &AuctionError{Message: "invalid toTokenAmount"}
		}

		coefficient := amount.Uint64()
		if coefficient > 0xffff {
			return nil, &AuctionError{Message: "coefficient exceeds uint16"}
		}
		coefficientBytes := make([]byte, 2)
		binary.BigEndian.PutUint16(coefficientBytes, uint16(coefficient))
		result = append(result, coefficientBytes...)
	}

	return result, nil
}
