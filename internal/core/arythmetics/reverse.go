package arythmetics

import (
	"fmt"
	"math/big"

	"github.com/paulpaulych/crypto/internal/core/rand"
)

const (
	maxTryCount = 100
)

// Reverse returns d=(rand^(-1))): rand*d(mod m) = 1
func Reverse(c, mod *big.Int) (*big.Int, error) {
	gcdExOut := GcdEx(mod, c)

	if gcdExOut.gcd.Cmp(big.NewInt(1)) != 0 {
		return nil, fmt.Errorf("can't find reverse: %v and %v aren't coprime", c, mod)
	}

	reverse := gcdExOut.y

	if reverse.Cmp(big.NewInt(0)) == -1 {
		return reverse.Add(reverse, mod), nil
	}
	return reverse, nil
}

// RandWithReverse for given 'P' returns rand,d satisfying
//  1. 1 < rand < P
//  2. rand*d = 1 (mod P)
func RandWithReverse(P *big.Int, random rand.Random) (r, d *big.Int, e error) {
	var reverse *big.Int
	predicate := func(c *big.Int) bool {
		d, e = Reverse(c, P)
		if e != nil {
			return false
		} else {
			reverse = d
			return true
		}
	}
	condRandom := rand.ConditionalRandom(
		maxTryCount,
		predicate,
		rand.FromToRandom(big.NewInt(2), P, random),
	)
	value, e := condRandom()
	if e != nil {
		return nil, nil, e
	}
	return value, reverse, nil
}

func CoprimeToRand(a *big.Int, r rand.Random) rand.Random {
	hasReverseElement := func(c *big.Int) bool {
		_, e := Reverse(c, a)
		return e != nil
	}
	return rand.ConditionalRandom(
		hasReverseElement,
		rand.FromToRandom(big.NewInt(2), a, r),
	)
}
