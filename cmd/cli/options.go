package cli

import (
	"fmt"
	"math/big"
)

type BigIntOpt struct {
	Value *big.Int
}

func (i *BigIntOpt) UnmarshalFlag(value string) error {
	intValue, success := new(big.Int).SetString(value, 10)
	if !success {
		return fmt.Errorf("cannot parse '%v' as integer value", value)
	}
	i.Value = intValue
	return nil
}

func (i BigIntOpt) MarshalFlag() (string, error) {
	return fmt.Sprintf("%v", i), nil
}
