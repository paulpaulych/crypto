package elgamal_cipher

import (
	"fmt"
	. "github.com/paulpaulych/crypto/internal/app/algorithms/diffie-hellman"
	"log"
	. "math/big"
	"reflect"
	"testing"
)

type testCase = struct {
	p   *Int
	bob *Bob
	msg *Int
}

func TestShamir(t *testing.T) {
	tests := append(
		generateCasesForAAndB(&Bob{
			CommonPub: &CommonPublicKey{P: NewInt(5), G: NewInt(3)},
			Pub:       NewInt(4),
			sec:       NewInt(2),
		}),
		testCase{
			bob: &Bob{
				CommonPub: &CommonPublicKey{P: NewInt(30803), G: NewInt(2)},
				Pub:       NewInt(8),
				sec:       NewInt(3),
			},
			msg: NewInt(1234),
		},
		testCase{
			bob: &Bob{
				CommonPub: &CommonPublicKey{P: NewInt(30803), G: NewInt(2)},
				Pub:       NewInt(1024),
				sec:       NewInt(10),
			},
			msg: NewInt(15),
		},
		testCase{
			bob: &Bob{
				CommonPub: &CommonPublicKey{P: NewInt(30803), G: NewInt(2)},
				Pub:       NewInt(256),
				sec:       NewInt(8),
			},
			msg: NewInt(1),
		},
		testCase{
			bob: &Bob{
				CommonPub: &CommonPublicKey{P: NewInt(30803), G: NewInt(2)},
				Pub:       NewInt(28273),
				sec:       NewInt(14807),
			},
			msg: NewInt(50),
		},
	)
	for _, tt := range tests {
		testName := fmt.Sprintf("testCase for B={commonPub={P=%v,G=%v}, sec=%v, Pub=%v}, msg=%v",
			tt.bob.CommonPub.P, tt.bob.CommonPub.G, tt.bob.sec, tt.bob.Pub, tt.msg)

		t.Run(testName, func(t *testing.T) {
			alice := NewAlice(tt.bob.CommonPub, tt.bob.Pub)
			encoded := alice.Encode(tt.msg, func(max *Int) (*Int, error) {
				return NewInt(randomInt), nil
			})
			log.Printf("ELGAMAL: R=%v, E=%v", encoded.R, encoded.E)
			decoded := tt.bob.Decode(encoded)
			log.Printf("decoded: %v", decoded)
			if !reflect.DeepEqual(decoded, tt.msg) {
				t.Errorf("decoded = %v, want %v", decoded, tt.msg)
			}
		})
	}
}

func generateCasesForAAndB(bob *Bob) []testCase {
	total := new(Int).Sub(bob.CommonPub.P, NewInt(3)).Int64()
	cases := make([]testCase, total)
	for i := int64(0); i < total; i++ {
		cases[i] = testCase{
			bob: bob,
			msg: NewInt(i + 2),
		}
	}
	return cases
}
