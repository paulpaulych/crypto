package elgamal_cipher

import (
	dh "github.com/paulpaulych/crypto/internal/app/algorithms/diffie-hellman"
	. "math/big"
	"reflect"
	"testing"
)

const randomInt = 21353

func TestBob_Decode(t *testing.T) {
	tests := []struct {
		name    string
		bob     Bob
		encoded Encoded
		want    *Int
	}{
		{
			name: "1",
			bob: Bob{
				CommonPub: &dh.CommonPublicKey{P: NewInt(30803), G: NewInt(2)},
				Pub:       NewInt(28273),
				sec:       NewInt(14807),
			},
			encoded: Encoded{R: NewInt(8), E: NewInt(21353)},
			want:    NewInt(50),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := tt.bob
			if got := b.Decode(&tt.encoded); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Decode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAlice_Encode(t *testing.T) {
	tests := []struct {
		name  string
		alice Alice
		msg   *Int
		want  *Encoded
	}{
		{
			name: "2",
			alice: Alice{
				CommonPub: &dh.CommonPublicKey{P: NewInt(30803), G: NewInt(2)},
				BobPub:    NewInt(28273),
			},
			msg:  NewInt(50),
			want: &Encoded{R: NewInt(8), E: NewInt(21353)},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := tt.alice
			if got := a.Encode(tt.msg, func(max *Int) (*Int, error) {
				return NewInt(randomInt), nil
			}); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Encode() = %v, want %v", got, tt.want)
			}
		})
	}
}
