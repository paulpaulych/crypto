package shamir

import (
	"crypto/rand"
	"errors"
	"fmt"
	"github.com/paulpaulych/crypto/internal/app/algorithms/arythmetics"
	"log"
	. "math/big"
)

type Alice struct {
	P    *Int
	C, D *Int
}

func (a *Alice) Step1(msg *Int) (*Int, error) {
	if msg.Cmp(a.P) != -1 {
		err := fmt.Sprintf("Shamir FATAL: msg=%v cannot be greater than P=%v", msg, a.P)
		return nil, errors.New(err)
	}
	if msg.Cmp(NewInt(1)) != 1 {
		err := fmt.Sprintf("Shamir FATAL: msg=%v cannot be less than 2", msg)
		return nil, errors.New(err)
	}
	return arythmetics.PowByMod(msg, a.C, a.P), nil
}

func (a *Alice) Step3(step2Out *Int) *Int {
	return arythmetics.PowByMod(step2Out, a.D, a.P)
}

func InitAlice(p *Int) (*Alice, error) {
	randomInt := func(max *Int) *Int {
		v, _ := rand.Int(rand.Reader, max)
		return v
	}
	c, d, err := initNode(p, randomInt)
	if err != nil {
		errMsg := fmt.Sprintf("writing step3out failed: %v", err)
		return nil, errors.New(errMsg)
	}
	alice := &Alice{P: p, C: c, D: d}
	log.Printf("Alice initialized: {P=%v,c=%v,d=%v}", alice.P, alice.C, alice.D)
	return alice, nil
}
