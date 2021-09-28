package shamir

import (
	"github.com/paulpaulych/crypto/commons"
	"log"
	. "math/big"
)

// initNode for given P returns c,d satisfying (c*d)mod(p-1) = 1
func initNode(p *Int, rand func(max *Int) *Int) (c, d *Int, e error) {
	max := new(Int).Sub(p, NewInt(2))

	for {
		rFrom0ToMax := rand(new(Int).Sub(max, NewInt(1)))
		c = new(Int).Add(rFrom0ToMax, NewInt(2))
		d, e = commons.Reverse(c, new(Int).Sub(p, NewInt(1)))

		if e != nil {
			log.Printf("shamir node initialization failed: '%s'. Retrying...", e)
			continue
		}

		return c, d, nil
	}
}
