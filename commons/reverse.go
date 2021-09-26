package commons

import "math/big"

// Reverse returns d=(c^(-1))): c*d(mod m) = 1
func Reverse(c, mod *big.Int) *big.Int {
	d := GcdEx(mod, c).y

	if d.Cmp(big.NewInt(0)) == -1 {
		return d.Add(d, mod)
	}
	return d
}
