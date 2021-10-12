package diffie_hellman

import (
	"github.com/paulpaulych/crypto/internal/app/algorithms/arythmetics"
	"math/big"
)

type Node struct {
	p, g, secretKey *big.Int
}

func (node *Node) CalcCommonKey(anotherPubKey *big.Int) CommonKey {
	return CommonKey{arythmetics.PowByMod(anotherPubKey, node.secretKey, node.p)}
}

type CommonKey = struct {
	value *big.Int
}
