package diffie_hellman

import (
	"github.com/paulpaulych/crypto/commons"
	"math/big"
)

type Node struct {
	p, g, secretKey *big.Int
}

func (node *Node) CalcCommonKey(anotherPubKey *big.Int) CommonKey {
	return CommonKey{commons.PowByMod(anotherPubKey, node.secretKey, node.p)}
}

type CommonKey = struct {
	value *big.Int
}
