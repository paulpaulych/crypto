package arythmetics

import (
	"log"
	. "math/big"
	"reflect"
	"testing"
)

func TestPowByMod(t *testing.T) {
	type args struct {
		x   *Int
		pow *Int
		mod *Int
	}
	tests := []struct {
		name string
		args args
		want *Int
	}{
		{
			name: "1 ^ 2 mod 10",
			args: args{NewInt(1), NewInt(2), NewInt(10)},
			want: NewInt(1),
		},
		{
			name: "2 ^ 8 mod 10",
			args: args{NewInt(2), NewInt(8), NewInt(10)},
			want: NewInt(6),
		},
		{
			name: "3 ^ 7 mod 10",
			args: args{NewInt(3), NewInt(7), NewInt(10)},
			want: NewInt(7),
		},
		{
			name: "7 ^ 19 mod 100",
			args: args{NewInt(7), NewInt(19), NewInt(100)},
			want: NewInt(43),
		},
		{
			name: "7 ^ 57 mod 100",
			args: args{NewInt(7), NewInt(57), NewInt(100)},
			want: NewInt(7),
		},
		{
			name: "3^35759111 mod 123654861",
			args: args{s2BigInt("3"), s2BigInt("35759111"), s2BigInt("123654861")},
			want: s2BigInt("32845410"),
		},
		{
			name: "102525021^91520007 mod 123654861",
			args: args{s2BigInt("102525021"), s2BigInt("91520007"), s2BigInt("123654861")},
			want: s2BigInt("32845410"),
		},
		{
			name: "102525021^91520007 mod 123654861",
			args: args{s2BigInt("102525021"), s2BigInt("91520007"), s2BigInt("123654861")},
			want: s2BigInt("32845410"),
		},
		{
			name: "102525021^91520007 mod 123654861",
			args: args{s2BigInt("102525021"), s2BigInt("91520007"), s2BigInt("123654861")},
			want: s2BigInt("32845410"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := PowByMod(tt.args.x, tt.args.pow, tt.args.mod)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PowByMod() = %v, want %v", got, tt.want)
			}
		})
	}
}

func s2BigInt(s string) *Int {
	v, success := new(Int).SetString(s, 10)
	if success != true {
		log.Fatalf("cannot convert '%s' to Int", s)
	}
	return v
}
