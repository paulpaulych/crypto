package commons

import (
	"errors"
	. "math/big"
	"reflect"
	"testing"
)

func TestReverse(t *testing.T) {
	type args struct {
		c   int64
		mod int64
	}
	type res struct {
		v *Int
		e error
	}
	tests := []struct {
		name string
		args args
		want res
	}{
		{"for 3 mod 7", args{3, 7}, res{NewInt(5), nil}},
		{"for 5 mod 8", args{5, 8}, res{NewInt(5), nil}},
		{"for 3 mod 53", args{3, 53}, res{NewInt(18), nil}},
		{"for 10 mod 53", args{10, 53}, res{NewInt(16), nil}},
		{"for 12 mod 16", args{12, 16}, res{nil, errors.New("can't find reverse: 12 and 16 aren't mutually simple")}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotV, gotE := Reverse(NewInt(tt.args.c), NewInt(tt.args.mod))
			if !reflect.DeepEqual(gotV, tt.want.v) {
				t.Errorf("Reverse() = %v, want %v", gotV, tt.want)
			}
			if gotE != tt.want.e && gotE.Error() != tt.want.e.Error() {
				t.Errorf("Reverse() = '%s', want '%s'", gotE, tt.want.e)
			}
		})
	}
}
