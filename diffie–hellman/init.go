package diffie_hellman

import (
	"crypto/rand"
	"math/big"
)

func InitDiffieHellmanNode(primeBits int) (Node, error) {
	var p, g *big.Int

	p, err := rand.Prime(rand.Reader, primeBits)
	if err != nil {
		return Node{}, err
	}
	g, err = rand.Int(rand.Reader, g.Sub(p, big.NewInt(3)))
	if err != nil {
		return Node{}, err
	}

	secretKey, err := generateSecretKey(p)
	if err != nil {
		return Node{}, err
	}

	return Node{
		secretKey: secretKey,
		p:         p,
		g:         g,
	}, nil
}

func generateSecretKey(p *big.Int) (*big.Int, error) {
	return rand.Int(rand.Reader, p)
}
