package core

import (
	"github.com/paulpaulych/crypto/commons"
	"log"
	"math/big"
)

type Alice struct {
	p    *big.Int
	c, d *big.Int
}

func (a *Alice) Step1(msg *big.Int) *big.Int {
	if msg.Cmp(a.p) != -1 {
		log.Fatalf("Shamir FATAL: msg=%v cannot be greater than p=%v", msg, a.p)
	}
	return commons.PowByMod(msg, a.c, a.p)
}

func (a *Alice) Step3(step2Out *big.Int) *big.Int {
	return commons.PowByMod(step2Out, a.d, a.p)
}
