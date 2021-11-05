package rsa_ds

import (
	"fmt"
	"github.com/paulpaulych/crypto/internal/core/arythmetics"
	"github.com/paulpaulych/crypto/internal/core/rand"
	. "math/big"
)

type PubKey struct {
	N   *Int
	Exp *Int
}

type SecretKey struct {
	N   *Int
	Exp *Int
}

type Signature struct {
	Value *Int
}

type Msg struct {
	Value *Int
}

func GenKeys(p, q *Int, random rand.Random) (*PubKey, *SecretKey, error) {
	N := new(Int).Mul(p, q)
	fi := new(Int).Mul(
		new(Int).Sub(p, NewInt(1)),
		new(Int).Sub(q, NewInt(1)),
	)
	c, d, err := arythmetics.RandWithReverse(fi, random)
	if err != nil {
		return nil, nil, err
	}

	return &PubKey{N: N, Exp: d}, &SecretKey{Exp: c, N: N}, nil
}

func Sign(key *SecretKey, msg *Msg, hashFn HashFn) (*Signature, error) {
	hash, err := hashFn(msg.Value)
	if err != nil {
		return nil, fmt.Errorf("can't stupidHash: %v", err)
	}
	return &Signature{
		Value: arythmetics.PowByMod(hash, key.Exp, key.N),
	}, err
}

func IsSignatureValid(key *PubKey, msg *Msg, sign *Signature, hashFn HashFn) (bool, error) {
	expectedHash, err := hashFn(msg.Value)
	if err != nil {
		return false, fmt.Errorf("can't stupidHash: %v", err)
	}
	hash := arythmetics.PowByMod(sign.Value, key.Exp, key.N)
	if expectedHash.Cmp(hash) != 0 {
		return false, nil
	}
	return true, nil
}

type HashFn func(orig *Int) (*Int, error)
