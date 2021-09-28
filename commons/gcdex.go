package commons

import (
	"log"
	"math/big"
)

type GcdExRes = struct {
	gcd  *big.Int
	x, y *big.Int
}

// GcdEx returns
//	gcd - GCD(a, b)
//  x, y corresponding ax+by=gcd
// Require: a > b
func GcdEx(a, b *big.Int) GcdExRes {
	switch a.Cmp(b) {
	case -1:
		log.Fatalf("GcdEx: A=%v cannot be less that B=%v", a, b)
	case 0:
		log.Fatalf("GcdEx: A=%v cannot be equal to B=%v", a, b)
	}
	if b.Cmp(big.NewInt(1)) == -1 {
		log.Fatalf("GcdEx: B=%v must be positive", b)
	}

	type row = struct{ i1, i2, i3 *big.Int }
	U := &row{new(big.Int).Set(a), big.NewInt(1), big.NewInt(0)}
	V := &row{new(big.Int).Set(b), big.NewInt(0), big.NewInt(1)}

	q := new(big.Int)
	for V.i1.Cmp(big.NewInt(0)) != 0 {
		q.Div(U.i1, V.i1)

		t1, t2, t3 := new(big.Int), new(big.Int), new(big.Int)

		t1.Mod(U.i1, V.i1)
		t2.Mul(V.i2, q).Sub(U.i2, t2)
		t3.Mul(V.i3, q).Sub(U.i3, t3)

		*U.i1 = *V.i1
		*U.i2 = *V.i2
		*U.i3 = *V.i3

		*V.i1 = *t1
		*V.i2 = *t2
		*V.i3 = *t3
	}
	return GcdExRes{gcd: U.i1, x: U.i2, y: U.i3}
}
