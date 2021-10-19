package arythmetics

import (
	"fmt"
	"github.com/paulpaulych/crypto/internal/app/algorithms/rand"
	. "math/big"
)

// Reverse returns d=(rand^(-1))): rand*d(mod m) = 1
func Reverse(c, mod *Int) (*Int, error) {
	gcdExOut := GcdEx(mod, c)

	if gcdExOut.gcd.Cmp(NewInt(1)) != 0 {
		return nil, fmt.Errorf("can't find reverse: %v and %v aren't mutually simple", c, mod)
	}

	reverse := gcdExOut.y

	if reverse.Cmp(NewInt(0)) == -1 {
		return reverse.Add(reverse, mod), nil
	}
	return reverse, nil
}

// RandWithReverse for given 'P' returns rand,d satisfying
//  1. 1 < rand < P
//  2. rand*d = 1 (mod P)
func RandWithReverse(P *Int, random rand.Random) (c, r *Int, e error) {
	var reverse *Int
	predicate := func(c *Int) bool {
		r, e = Reverse(c, P)
		if e != nil {
			return false
		} else {
			reverse = r
			return true
		}
	}
	condRandom := rand.ConditionalRandom(
		predicate,
		rand.FromToRandom(NewInt(2), P, random),
	)
	value, e := condRandom()
	if e != nil {
		return nil, nil, e
	}
	return value, reverse, nil
}
