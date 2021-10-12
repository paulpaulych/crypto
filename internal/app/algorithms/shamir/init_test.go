package shamir

import (
	. "math/big"
	"testing"
)

func Test_initNode(t *testing.T) {
	tests := []struct {
		name string
		p    *Int
		c    *Int
	}{
		{"case1", NewInt(17), NewInt(3)},
		{"case1", NewInt(3), NewInt(2)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotC, gotD, err := initNode(tt.p, func(max *Int) *Int { return tt.c })
			if err != nil {
				t.Errorf("initNode() error = %v", err)
				return
			}
			res := new(Int)
			if res.Mul(gotC, gotD).Mod(res, new(Int).Sub(tt.p, NewInt(1))).Cmp(NewInt(1)) != 0 {
				t.Errorf("gotD=%v is not reverse for gotC=%v", gotD, gotC)
			}
		})
	}
}
