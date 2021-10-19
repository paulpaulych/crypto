package rsa_cipher

import (
	. "github.com/paulpaulych/crypto/internal/app/algorithms/arythmetics"
	. "github.com/paulpaulych/crypto/internal/app/algorithms/rand"
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
	c, d, err := RandWithReverse(fi, random)
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
