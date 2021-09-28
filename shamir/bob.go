package shamir

import (
	"github.com/paulpaulych/crypto/commons"
	"math/big"
)

type Bob struct {
	p    *big.Int
	c, d *big.Int
}

func (b *Bob) Step2(step1Out *big.Int) *big.Int {
	return commons.PowByMod(step1Out, b.c, b.p)
}

func (b *Bob) Decode(step3Out *big.Int) *big.Int {
	return commons.PowByMod(step3Out, b.d, b.p)
}
