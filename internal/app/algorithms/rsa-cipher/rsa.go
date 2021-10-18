package rsa_cipher

import (
	. "github.com/paulpaulych/crypto/internal/app/algorithms/arythmetics"
	. "github.com/paulpaulych/crypto/internal/app/algorithms/rand"
	"log"
	. "math/big"
)

type BobSecret struct {
	p, q *Int
	c    *Int
}
type BobPub struct {
	N *Int
	D *Int
}

type Bob struct {
	*BobSecret
	*BobPub
}

type Encoded struct {
	Value *Int
}

type Alice struct {
	BobPub *BobPub
}

func NewAlice(bobPub *BobPub) *Alice {
	return &Alice{BobPub: bobPub}
}

func (a Alice) Encode(msg *Int) *Encoded {
	return &Encoded{
		Value: PowByMod(msg, a.BobPub.D, a.BobPub.N),
	}
}

func NewBob(p *Int, q *Int, random Random) (*Bob, error) {
	N := new(Int).Mul(p, q)
	fi := new(Int).Mul(
		new(Int).Sub(p, NewInt(1)),
		new(Int).Sub(q, NewInt(1)),
	)
	c, d, err := initNode(fi, random)
	if err != nil {
		return nil, err
	}

	return &Bob{
		BobSecret: &BobSecret{p: p, q: q, c: c},
		BobPub:    &BobPub{N: N, D: d},
	}, nil
}

func (b Bob) Decode(encoded *Encoded) *Int {
	return PowByMod(encoded.Value, b.BobSecret.c, b.BobPub.N)
}

// initNode for given fi returns c,D satisfying (c*D) (mod fi) = 1
func initNode(fi *Int, rand Random) (c, d *Int, e error) {
	fromToRandom := FromToRandom(NewInt(2), fi, rand)
	for {
		c, e := fromToRandom()
		if e != nil {
			return nil, nil, e
		}
		d, e = Reverse(c, fi)

		if e != nil {
			log.Printf("shamir-cipher node initialization failed: %s. Retrying...", e)
			continue
		}

		return c, d, nil
	}
}
