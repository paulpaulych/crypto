package commons

import (
	"math/big"
	"reflect"
	"testing"
)

func TestPowByMod(t *testing.T) {
	type args struct {
		x   int64
		pow int64
		mod int64
	}
	tests := []struct {
		name string
		args args
		want int64
	}{
		{
			name: "1 ^ 2 mod 10",
			args: args{1, 2, 10},
			want: 1,
		},
		{
			name: "2 ^ 8 mod 10",
			args: args{2, 8, 10},
			want: 6,
		},
		{
			name: "3 ^ 7 mod 10",
			args: args{3, 7, 10},
			want: 7,
		},
		{
			name: "7 ^ 19 mod 100",
			args: args{7, 19, 100},
			want: 43,
		},
		{
			name: "7 ^ 57 mod 100",
			args: args{7, 57, 100},
			want: 7,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := PowByMod(big.NewInt(tt.args.x), big.NewInt(tt.args.pow), big.NewInt(tt.args.mod))
			if !reflect.DeepEqual(got, big.NewInt(tt.want)) {
				t.Errorf("PowByMod(%x) = %x, want %x", tt.args, got, tt.want)
			}
		})
	}
}
