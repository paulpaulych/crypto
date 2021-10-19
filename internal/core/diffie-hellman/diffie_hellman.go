package diffie_hellman

import (
	"fmt"
	. "math/big"
)

const minP = 256

type CommonPublicKey interface {
	P() *Int
	G() *Int
}

type commonPublicKey struct {
	p, g *Int
}

func (k commonPublicKey) P() *Int {
	return k.p
}

func (k commonPublicKey) G() *Int {
	return k.g
}

func NewCommonPublicKey(p, g *Int) (CommonPublicKey, error) {
	if p.Cmp(NewInt(minP)) < 0 {
		return nil, fmt.Errorf("p=%v cannot be less than %v", p, minP)
	}
	//	TODO: check that G is generator of multiplicative group modulo P if possible
	return commonPublicKey{p: p, g: g}, nil
}
