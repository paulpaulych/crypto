package arythmetics

import (
	"fmt"
	"math/big"
)

// Reverse returns d=(c^(-1))): c*d(mod m) = 1
func Reverse(c, mod *big.Int) (*big.Int, error) {
	gcdExOut := GcdEx(mod, c)

	if gcdExOut.gcd.Cmp(big.NewInt(1)) != 0 {
		return nil, fmt.Errorf("can't find reverse: %v and %v aren't mutually simple", c, mod)
	}

	reverse := gcdExOut.y

	if reverse.Cmp(big.NewInt(0)) == -1 {
		return reverse.Add(reverse, mod), nil
	}
	return reverse, nil
}
