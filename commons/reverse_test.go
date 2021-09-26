package commons

import (
	"math/big"
	"reflect"
	"testing"
)

func TestReverse(t *testing.T) {
	type args struct {
		c   int64
		mod int64
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{"for 3 mod 7", args{3, 7}, 5},
		{"for 5 mod 8", args{5, 8}, 5},
		{"for 3 mod 53", args{3, 53}, 18},
		{"for 10 mod 53", args{10, 53}, 16},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Reverse(big.NewInt(tt.args.c), big.NewInt(tt.args.mod))
			if !reflect.DeepEqual(got, big.NewInt(tt.want)) {
				t.Errorf("Reverse() = %v, want %v", got, tt.want)
			}
		})
	}
}
