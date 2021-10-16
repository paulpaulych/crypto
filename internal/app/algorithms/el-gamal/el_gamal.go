package shamir_cipher

import (
	"crypto/rand"
	. "github.com/paulpaulych/crypto/internal/app/algorithms/arythmetics"
	dh "github.com/paulpaulych/crypto/internal/app/algorithms/diffie-hellman"
	. "math/big"
)

type Encoded struct {
	r, e *Int
}

func Send(cPub *dh.CommonPublicKey, bobPub *Int, msg *Int) *Encoded {
	secret, err := rand.Int(rand.Reader, cPub.P)
	if err != nil {
		return nil
	}

	tmp := new(Int)
	tmp.Mul(msg, PowByMod(bobPub, secret, cPub.P))
	e := tmp.Mod(tmp, cPub.P)
	return &Encoded{
		r: PowByMod(cPub.G, secret, cPub.P),
		e: e,
	}
}

type Bob struct {
	CommonPub dh.CommonPublicKey
	Pub       *Int
	sec       *Int
}

func (b Bob) Recv(encoded *Encoded) *Int {
	p := b.CommonPub.P
	tmp := new(Int)
	tmp.Sub(p, NewInt(1))
	tmp.Sub(tmp, b.sec)
	return tmp.
		Mul(encoded.e, PowByMod(encoded.r, tmp, p)).
		Mod(tmp, p)
}
