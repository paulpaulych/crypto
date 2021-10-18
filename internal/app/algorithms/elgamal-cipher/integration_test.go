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
	commonPubKey, e := NewCommonPublicKey(NewInt(30803), NewInt(2))
	if e != nil {
		log.Panicf("failed to initialize test data: %v", e)
	}

	tests := []testCase{
		{
			bob: &Bob{
				CommonPub: commonPubKey,
				Pub:       NewInt(8),
				sec:       NewInt(3),
			},
			msg: NewInt(1234),
		},
		{
			bob: &Bob{
				CommonPub: commonPubKey,
				Pub:       NewInt(1024),
				sec:       NewInt(10),
			},
			msg: NewInt(15),
		},
		{
			bob: &Bob{
				CommonPub: commonPubKey,
				Pub:       NewInt(256),
				sec:       NewInt(8),
			},
			msg: NewInt(1),
		},
		{
			bob: &Bob{
				CommonPub: commonPubKey,
				Pub:       NewInt(28273),
				sec:       NewInt(14807),
			},
			msg: NewInt(50),
		},
	}
	for _, tt := range tests {
		testName := fmt.Sprintf("testCase for B={commonPubKey={P=%v,G=%v}, sec=%v, Pub=%v}, msg=%v",
			tt.bob.CommonPub.P(), tt.bob.CommonPub.G(), tt.bob.sec, tt.bob.Pub, tt.msg)

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
