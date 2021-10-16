package diffie_hellman

import (
	"math/big"
)

type CommonPublicKey struct {
	P, G *big.Int
}
