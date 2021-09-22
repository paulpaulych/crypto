package commons

import (
	"fmt"
	"math/big"
	"reflect"
	"testing"
)

func TestGcd(t *testing.T) {
	type args struct {
		a, b int64
	}
	tests := []struct {
		args args
		want int64
	}{
		{args{1, 2}, 1},
		{args{10, 30}, 10},
		{args{10, 35}, 5},
		{args{51, 34}, 17},
		{args{18, 18}, 18},
		{args{1, 1}, 1},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case%v", i), func(t *testing.T) {
			got := Gcd(big.NewInt(tt.args.a), big.NewInt(tt.args.b))
			if !reflect.DeepEqual(got, big.NewInt(tt.want)) {
				t.Errorf("Gcd() = %v, want %v", got, tt.want)
			}
		})
	}
}
