package rsa_ds

import (
	"fmt"
	"github.com/paulpaulych/crypto/internal/app/algorithms/arythmetics"
	"github.com/paulpaulych/crypto/internal/app/algorithms/rand"
	. "math/big"
)

type PubKey struct {
	N   *Int
	Exp *Int
}

type SecKey struct {
	N   *Int
	Exp *Int
}

type Signed struct {
	Msg       *Int
	Signature *Int
}

func GenKeys(p, q *Int, random rand.Random) (*PubKey, *SecKey, error) {
	N := new(Int).Mul(p, q)
	fi := new(Int).Mul(
		new(Int).Sub(p, NewInt(1)),
		new(Int).Sub(q, NewInt(1)),
	)
	c, d, err := arythmetics.RandWithReverse(fi, random)
	if err != nil {
		return nil, nil, err
	}

	return &PubKey{N: N, Exp: d}, &SecKey{Exp: c}, nil
}

func Sign(key *SecKey, msg *Int, hashFn HashFn) (*Signed, error) {
	hash, err := hashFn(msg)
	if err != nil {
		return nil, fmt.Errorf("can't stupidHash: %v", err)
	}
	signature := arythmetics.PowByMod(hash, key.Exp, key.N)
	return &Signed{
		Msg:       msg,
		Signature: signature,
	}, err
}

func IsSignatureValid(key *PubKey, signed *Signed, hashFn HashFn) (bool, error) {
	expectedHash, err := hashFn(signed.Msg)
	if err != nil {
		return false, fmt.Errorf("can't stupidHash: %v", err)
	}
	hash := arythmetics.PowByMod(signed.Signature, key.Exp, key.N)
	if expectedHash.Cmp(hash) != 0 {
		return false, nil
	}
	return true, nil
}

type HashFn func(orig *Int) (*Int, error)
