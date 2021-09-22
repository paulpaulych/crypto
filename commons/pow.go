package commons

import "math/big"

func PowByMod(x, pow, mod *big.Int) *big.Int {
	res := big.NewInt(1)
	ai := x
	iterateBin(pow, func(isOne bool) {
		if isOne {
			res.Mul(res, ai).Mod(res, mod)
		}
		pow2(ai, mod)
	})
	return res
}

func pow2(x, mod *big.Int) *big.Int {
	return x.Mul(x, x).Mod(x, mod)
}

var (
	zero = big.NewInt(0)
	two  = big.NewInt(2)
)

func iterateBin(x *big.Int, action func(isOne bool)) {
	shift := uint(0)
	ti := new(big.Int)
	for ; ti.Rsh(x, shift).Cmp(zero) != 0; shift++ {
		mod := new(big.Int)
		action(mod.Mod(ti, two).Cmp(zero) != 0)
	}
}
