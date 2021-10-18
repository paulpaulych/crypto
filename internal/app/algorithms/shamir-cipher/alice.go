package shamir_cipher

import (
	"fmt"
	"github.com/paulpaulych/crypto/internal/app/algorithms/arythmetics"
	"github.com/paulpaulych/crypto/internal/app/algorithms/rand"
	"log"
	. "math/big"
)

type Alice struct {
	P    *Int
	C, D *Int
}

func (a *Alice) Step1(msg *Int) (*Int, error) {
	if msg.Cmp(a.P) != -1 {
		return nil, fmt.Errorf("shamir-cipher FATAL: msg=%v cannot be greater than P=%v", msg, a.P)
	}
	if msg.Cmp(NewInt(1)) != 1 {
		return nil, fmt.Errorf("Shamir FATAL: msg=%v cannot be less than 2", msg)
	}
	return arythmetics.PowByMod(msg, a.C, a.P), nil
}

func (a *Alice) Step3(step2Out *Int) *Int {
	return arythmetics.PowByMod(step2Out, a.D, a.P)
}

func InitAlice(p *Int) (*Alice, error) {
	c, d, err := initNode(p, rand.CryptoSafeRandom())
	if err != nil {
		return nil, fmt.Errorf("writing step3out failed: %v", err)
	}
	alice := &Alice{P: p, C: c, D: d}
	log.Printf("Alice initialized: {P=%v,c=%v,d=%v}", alice.P, alice.C, alice.D)
	return alice, nil
}
