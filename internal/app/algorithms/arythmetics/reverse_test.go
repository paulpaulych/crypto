package arythmetics

import (
	"errors"
	"github.com/paulpaulych/crypto/internal/app/algorithms/rand"
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

func TestRandWithReverse(t *testing.T) {
	tests := []struct {
		name  string
		P     *Int
		rand  rand.Random
		wantC *Int
		wantD *Int
	}{
		{
			name:  "1",
			P:     NewInt(17),
			rand:  rand.ConstRand(NewInt(3)),
			wantC: NewInt(5),
			wantD: NewInt(7),
		},
		{
			name:  "2",
			P:     NewInt(10),
			rand:  rand.CyclicRandom(NewInt(5), NewInt(3)),
			wantC: NewInt(7),
			wantD: NewInt(3),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotC, gotD, _ := RandWithReverse(tt.P, tt.rand)
			if gotC.Cmp(tt.wantC) != 0 {
				t.Errorf("gotC=%v , want %v", gotC, tt.wantC)
			}

			if gotD.Cmp(tt.wantD) != 0 {
				t.Errorf("gotD=%v , want %v", gotD, tt.wantD)
			}
		})
	}
}
