package diffie_hellman

import (
	"math/big"
	"reflect"
	"testing"
)

func TestNode_CalcCommonKey(t *testing.T) {
	type testNode struct{ p, g, secretKey int64 }
	tests := []struct {
		name         string
		fields       testNode
		srcPublicKey int64
		want         int64
	}{
		{"case1", testNode{7, 2, 3}, 3, 6},
		{"case2", testNode{7, 3, 2}, 3, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			node := &Node{
				p:         big.NewInt(tt.fields.p),
				g:         big.NewInt(tt.fields.g),
				secretKey: big.NewInt(tt.fields.secretKey),
			}
			got := node.CalcCommonKey(big.NewInt(tt.srcPublicKey))
			if !reflect.DeepEqual(got.value, big.NewInt(tt.want)) {
				t.Errorf("CalcCommonKey() = %v, want %v", got.value, tt.want)
			}
		})
	}
}
