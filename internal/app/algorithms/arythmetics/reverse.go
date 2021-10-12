package arythmetics

import (
	"errors"
	"fmt"
	"math/big"
)

// Reverse returns d=(c^(-1))): c*d(mod m) = 1
func Reverse(c, mod *big.Int) (*big.Int, error) {
	gcdExOut := GcdEx(mod, c)

	if gcdExOut.gcd.Cmp(big.NewInt(1)) != 0 {
		msg := fmt.Sprintf("can't find reverse: %v and %v aren't mutually simple", c, mod)
		return nil, errors.New(msg)
	}

	reverse := gcdExOut.y

	if reverse.Cmp(big.NewInt(0)) == -1 {
		return reverse.Add(reverse, mod), nil
	}
	return reverse, nil
}
