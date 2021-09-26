package commons

import (
	"fmt"
	"math/big"
	"testing"
)

func TestGcdEx(t *testing.T) {
	var assertEqual = func(got, want GcdExRes) {
		if got.gcd.Cmp(want.gcd) != 0 || got.x.Cmp(want.x) != 0 || got.y.Cmp(want.y) != 0 {
			t.Errorf("GcdEx() = %v, want %v", fmtRes(got), fmtRes(want))
		}
	}

	type args = struct{ a, b int64 }
	type res = struct{ gcd, x, y int64 }
	tests := []struct {
		name string
		args args
		want res
	}{
		{"for 3, 7", args{7, 3}, res{1, 1, -2}},
		{"for 3, 1", args{3, 1}, res{1, 0, 1}},
		{"for 18, 12", args{18, 12}, res{6, 1, -1}},
		{"for 28, 19", args{28, 19}, res{1, -2, 3}},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf(tt.name), func(t *testing.T) {
			got := GcdEx(big.NewInt(tt.args.a), big.NewInt(tt.args.b))
			want := GcdExRes{big.NewInt(tt.want.gcd), big.NewInt(tt.want.x), big.NewInt(tt.want.y)}
			assertEqual(got, want)
		})
	}
}

func fmtRes(v GcdExRes) string {
	return fmt.Sprintf("{gcd:%v, x=%v, y=%v}", v.gcd, v.x, v.y)
}
