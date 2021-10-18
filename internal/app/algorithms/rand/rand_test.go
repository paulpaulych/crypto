package rand

import (
	"fmt"
	. "math/big"
	"reflect"
	"testing"
)

func TestFromToRandom(t *testing.T) {
	tests := []struct {
		from *Int
		to   *Int
		rand Random
		want *Int
	}{
		{
			from: NewInt(20),
			to:   NewInt(30),
			rand: ConstRand(NewInt(20)),
			want: NewInt(20),
		},
		{
			from: NewInt(20),
			to:   NewInt(30),
			rand: ConstRand(NewInt(19)),
			want: NewInt(29),
		},
		{
			from: NewInt(20),
			to:   NewInt(30),
			rand: ConstRand(NewInt(18)),
			want: NewInt(28),
		},
		{
			from: NewInt(20),
			to:   NewInt(30),
			rand: ConstRand(NewInt(22)),
			want: NewInt(22),
		},
		{
			from: NewInt(20),
			to:   NewInt(30),
			rand: ConstRand(NewInt(30)),
			want: NewInt(20),
		},
		{
			from: NewInt(20),
			to:   NewInt(30),
			rand: ConstRand(NewInt(31)),
			want: NewInt(21),
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %v", i), func(t *testing.T) {
			rand := FromToRandom(tt.from, tt.to, tt.rand)
			if got, _ := rand(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FromToRandom() = %v, want %v", got, tt.want)
			}
		})
	}
}
