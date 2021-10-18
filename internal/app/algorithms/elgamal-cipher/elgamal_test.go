package elgamal_cipher

import (
	"fmt"
	dh "github.com/paulpaulych/crypto/internal/app/algorithms/diffie-hellman"
	"github.com/paulpaulych/crypto/internal/app/algorithms/rand"
	. "math/big"
	"reflect"
	"testing"
)

const randomInt = 1

var commonPub, _ = dh.NewCommonPublicKey(NewInt(30803), NewInt(2))
var bobPub = NewInt(28273)

func TestBob_Decode(t *testing.T) {
	tests := []struct {
		bob     Bob
		encoded Encoded
		want    *Int
	}{
		{
			bob: Bob{
				CommonPub: commonPub,
				Pub:       bobPub,
				sec:       NewInt(14807),
			},
			encoded: Encoded{R: NewInt(8), E: NewInt(21353)},
			want:    NewInt(50),
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %v", i), func(t *testing.T) {
			b := tt.bob
			if got := b.Decode(&tt.encoded); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Decode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAlice_Encode(t *testing.T) {
	tests := []struct {
		alice Alice
		msg   *Int
		want  *Encoded
	}{
		{
			alice: Alice{
				CommonPub: commonPub,
				BobPub:    bobPub,
			},
			msg:  NewInt(50),
			want: &Encoded{R: NewInt(8), E: NewInt(21353)},
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %v", i), func(t *testing.T) {
			a := tt.alice
			got := a.Encode(tt.msg, rand.ConstRand(NewInt(1)))
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encode() = %v, want %v", got, tt.want)
			}
		})
	}
}
