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

func TestCyclicRandom(t *testing.T) {
	tests := []struct {
		values []*Int
	}{
		{[]*Int{NewInt(1), NewInt(2), NewInt(3)}},
		{[]*Int{NewInt(1)}},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %v", i), func(t *testing.T) {
			rand := CyclicRandom(tt.values...)
			for i, v := range tt.values {
				if got, _ := rand(); !reflect.DeepEqual(got, v) {
					t.Errorf("rand() = %v, want %v on position %v", got, v, i)
				}
			}
		})
	}
}

func TestConditionalRandom(t *testing.T) {
	tests := []struct {
		name      string
		predicate func(*Int) bool
		rand      Random
		want      *Int
	}{
		{
			name: "find first even number equal to 2",
			rand: CyclicRandom(NewInt(1), NewInt(0), NewInt(2), NewInt(3)),
			predicate: func(i *Int) bool {
				if i.Cmp(NewInt(2)) == 0 {
					return true
				} else {
					return false
				}
			},
			want: NewInt(2),
		},
		{
			name: "returns number after second ZERO",
			rand: CyclicRandom(NewInt(1), NewInt(0), NewInt(2), NewInt(3)),
			predicate: func() func(i *Int) bool {
				zeroCnt := 0
				return func(i *Int) bool {
					if zeroCnt == 2 {
						return true
					}
					if i.Cmp(NewInt(0)) == 0 {
						zeroCnt++
					}
					return false
				}
			}(),
			want: NewInt(2),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rand := ConditionalRandom(tt.predicate, tt.rand)
			if got, _ := rand(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ConditionalRandom() = %v, want %v", got, tt.want)
			}
		})
	}
}
