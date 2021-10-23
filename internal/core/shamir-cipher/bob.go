package shamir_cipher

import (
	"github.com/paulpaulych/crypto/internal/core/arythmetics"
	"github.com/paulpaulych/crypto/internal/core/rand"
	"log"
	. "math/big"
)

type Bob struct {
	p    *Int
	c, d *Int
}

func (b *Bob) Step2(step1Out *Int) *Int {
	return arythmetics.PowByMod(step1Out, b.c, b.p)
}

func (b *Bob) Decode(step3Out *Int) *Int {
	return arythmetics.PowByMod(step3Out, b.d, b.p)
}

func InitBob(p *Int) (*Bob, error) {
	c, d, err := arythmetics.RandWithReverse(new(Int).Sub(p, NewInt(1)), rand.CryptoSafeRandom())
	if err != nil {
		return nil, err
	}
	bob := &Bob{p: p, c: c, d: d}
	log.Printf("Bob initialized: {P=%v,c=%v,d=%v}", bob.p, bob.c, bob.d)
	return bob, nil
}
