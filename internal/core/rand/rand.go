package rand

import (
	"crypto/rand"
	"fmt"
	"log"
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

// ConstRand is useful for tests
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

// ConditionalRandom returns first random value that matches the predicate
func ConditionalRandom(predicate func(*Int) bool, rand Random) Random {
	return func() (*Int, error) {
		for i := 0; ; i++ {
			value, e := rand()
			if e != nil {
				return nil, e
			}
			if !predicate(value) {
				log.Printf("ConditionalRandom: try %v failed. Retrying...", i)
				continue
			}
			return value, nil
		}
	}
}

// CyclicRandom is useful for tests
func CyclicRandom(values ...*Int) Random {
	if values == nil || len(values) == 0 {
		log.Fatalf("CyclicRandom: values cannot be empty")
	}

	i := 0
	size := len(values)
	return func() (*Int, error) {
		res := values[i]
		i = (i + 1) % size
		return res, nil
	}
}
