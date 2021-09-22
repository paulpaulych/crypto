package commons

import (
	"math/big"
)

func Gcd(a, b *big.Int) *big.Int {
	var x, y *big.Int
	if a.Cmp(b) > 0 {
		x, y = a, b
	} else {
		x, y = b, a
	}
	for rem := new(big.Int); y.Cmp(zero) != 0; {
		rem.Rem(x, y)
		*x = *y
		*y = *rem
	}
	return x
}
