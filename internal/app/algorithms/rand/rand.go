package rand

import (
	"crypto/rand"
	"fmt"
	. "math/big"
)

// TODO: increase max value?
var max, _ = new(Int).SetString("11111111111111111111111111111111111", 10)

type Random = func() (*Int, error)

func CryptoSafeRandom() Random {
	return func() (*Int, error) {
		return rand.Int(rand.Reader, max)
	}
}

func ConstRand(value *Int) Random {
	return func() (*Int, error) {
		return value, nil
	}
}

func FromToRandom(from *Int, to *Int, rand Random) Random {
	return func() (*Int, error) {
		if from.Cmp(to) >= 0 {
			return nil, fmt.Errorf("RANDOM: min=%v must be less than max=%v", from, to)
		}
		value, e := rand()
		if e != nil {
			return nil, e
		}
		diff := new(Int).Sub(to, from)
		shift := new(Int).Mod(value, diff)
		return new(Int).Add(shift, from), nil
	}
}
