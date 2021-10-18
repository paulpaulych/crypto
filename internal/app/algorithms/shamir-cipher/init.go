package shamir_cipher

import (
	"github.com/paulpaulych/crypto/internal/app/algorithms/arythmetics"
	"github.com/paulpaulych/crypto/internal/app/algorithms/rand"
	"log"
	. "math/big"
)

// initNode for given P returns c,d satisfying (c*d)mod(p-1) = 1
func initNode(p *Int, random rand.Random) (c, d *Int, e error) {
	pSub1 := new(Int).Sub(p, NewInt(1))
	fromToRand := rand.FromToRandom(
		NewInt(2),
		pSub1,
		random,
	)
	for {
		c, err := fromToRand()
		if err != nil {
			return nil, nil, err
		}

		d, e = arythmetics.Reverse(c, pSub1)

		if e != nil {
			log.Printf("shamir-cipher node initialization failed: %s. Retrying...", e)
			continue
		}

		return c, d, nil
	}
}
