package timelocks

import (
	"fmt"
	"math/big"
)

const (
	Web3Type                  = "uint256"
	DefaultRescueDelay        = 604800
	UINT_32_MAX        uint64 = 0xffffffff
)

type TimeLocks struct {
	srcWithdrawal         *big.Int
	srcPublicWithdrawal   *big.Int
	srcCancellation       *big.Int
	srcPublicCancellation *big.Int
	dstWithdrawal         *big.Int
	dstPublicWithdrawal   *big.Int
	dstCancellation       *big.Int
	deployedAt            *big.Int
}

func NewTimeLocks(params TimeLocksParams) (*TimeLocks, error) {
	if err := validateUint32(params.SrcWithdrawal); err != nil {
		return nil, fmt.Errorf("srcWithdrawal: %w", err)
	}
	if err := validateUint32(params.SrcPublicWithdrawal); err != nil {
		return nil, fmt.Errorf("srcPublicWithdrawal: %w", err)
	}
	if params.SrcWithdrawal.Cmp(params.SrcPublicWithdrawal) >= 0 {
		return nil, fmt.Errorf("srcWithdrawal must be < srcPublicWithdrawal")
	}
	if err := validateUint32(params.SrcCancellation); err != nil {
		return nil, fmt.Errorf("srcCancellation: %w", err)
	}
	if params.SrcPublicWithdrawal.Cmp(params.SrcCancellation) >= 0 {
		return nil, fmt.Errorf("srcPublicWithdrawal must be < srcCancellation")
	}
	if err := validateUint32(params.SrcPublicCancellation); err != nil {
		return nil, fmt.Errorf("srcPublicCancellation: %w", err)
	}
	if params.SrcCancellation.Cmp(params.SrcPublicCancellation) >= 0 {
		return nil, fmt.Errorf("srcCancellation must be < srcPublicCancellation")
	}
	if err := validateUint32(params.DstWithdrawal); err != nil {
		return nil, fmt.Errorf("dstWithdrawal: %w", err)
	}
	if params.DstWithdrawal.Cmp(params.DstPublicWithdrawal) >= 0 {
		return nil, fmt.Errorf("dstWithdrawal must be < dstPublicWithdrawal")
	}
	if err := validateUint32(params.DstPublicWithdrawal); err != nil {
		return nil, fmt.Errorf("dstPublicWithdrawal: %w", err)
	}
	if params.DstPublicWithdrawal.Cmp(params.DstCancellation) >= 0 {
		return nil, fmt.Errorf("dstPublicWithdrawal must be < dstCancellation")
	}
	if err := validateUint32(params.DstCancellation); err != nil {
		return nil, fmt.Errorf("dstCancellation: %w", err)
	}

	deployedAt := params.DeployedAt
	if deployedAt == nil {
		deployedAt = big.NewInt(0)
	}
	if err := validateUint32(deployedAt); err != nil {
		return nil, fmt.Errorf("deployedAt: %w", err)
	}

	return &TimeLocks{
		srcWithdrawal:         new(big.Int).Set(params.SrcWithdrawal),
		srcPublicWithdrawal:   new(big.Int).Set(params.SrcPublicWithdrawal),
		srcCancellation:       new(big.Int).Set(params.SrcCancellation),
		srcPublicCancellation: new(big.Int).Set(params.SrcPublicCancellation),
		dstWithdrawal:         new(big.Int).Set(params.DstWithdrawal),
		dstPublicWithdrawal:   new(big.Int).Set(params.DstPublicWithdrawal),
		dstCancellation:       new(big.Int).Set(params.DstCancellation),
		deployedAt:            new(big.Int).Set(deployedAt),
	}, nil
}

func TimeLocksFromBigInt(val *big.Int) (*TimeLocks, error) {
	params := make([]*big.Int, 8)
	mask := big.NewInt(0xffffffff)
	for i := 0; i < 8; i++ {
		shift := big.NewInt(int64(i * 32))
		shifted := new(big.Int).Rsh(val, uint(shift.Uint64()))
		params[i] = new(big.Int).And(shifted, mask)
	}

	return NewTimeLocks(TimeLocksParams{
		DeployedAt:            params[0],
		DstCancellation:       params[1],
		DstPublicWithdrawal:   params[2],
		DstWithdrawal:         params[3],
		SrcPublicCancellation: params[4],
		SrcCancellation:       params[5],
		SrcPublicWithdrawal:   params[6],
		SrcWithdrawal:         params[7],
	})
}

type TimeLocksParams struct {
	SrcWithdrawal         *big.Int
	SrcPublicWithdrawal   *big.Int
	SrcCancellation       *big.Int
	SrcPublicCancellation *big.Int
	DstWithdrawal         *big.Int
	DstPublicWithdrawal   *big.Int
	DstCancellation       *big.Int
	DeployedAt            *big.Int
}

func (tl *TimeLocks) Build() *big.Int {
	result := big.NewInt(0)
	values := []*big.Int{
		tl.deployedAt,
		tl.dstCancellation,
		tl.dstPublicWithdrawal,
		tl.dstWithdrawal,
		tl.srcPublicCancellation,
		tl.srcCancellation,
		tl.srcPublicWithdrawal,
		tl.srcWithdrawal,
	}
	for _, val := range values {
		result.Lsh(result, 32)
		result.Or(result, new(big.Int).And(val, big.NewInt(0xffffffff)))
	}
	return result
}

func (tl *TimeLocks) DeployedAt() *big.Int {
	return new(big.Int).Set(tl.deployedAt)
}

func (tl *TimeLocks) SetDeployedAt(time *big.Int) *TimeLocks {
	if err := validateUint32(time); err != nil {
		panic(err)
	}
	tl.deployedAt = new(big.Int).Set(time)
	return tl
}

func validateUint32(val *big.Int) error {
	if val.Sign() < 0 {
		return fmt.Errorf("value must be non-negative")
	}
	max := big.NewInt(int64(UINT_32_MAX))
	if val.Cmp(max) > 0 {
		return fmt.Errorf("value exceeds uint32 max")
	}
	return nil
}
