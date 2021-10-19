package rsa_cipher

import (
	"fmt"
	"github.com/paulpaulych/crypto/internal/core/rand"
	"log"
	. "math/big"
	"reflect"
	"testing"
)

type testCase = struct {
	bob *Bob
	msg *Int
}

func TestRsa(t *testing.T) {
	bob1, err := NewBob(NewInt(11), NewInt(13), rand.ConstRand(NewInt(5)))
	if err != nil {
		log.Panicf("cant init test data: %s", err)
	}
	tests := []testCase{
		{
			bob: bob1,
			msg: NewInt(13),
		},
	}
	for _, tt := range tests {
		testName := fmt.Sprintf("testCase for B={public={N=%v,D=%v}, sec={p=%v,q=%v,c=%v}, msg=%v",
			tt.bob.BobPub.N, tt.bob.BobPub.D, tt.bob.BobSecret.p, tt.bob.BobSecret.q, tt.bob.BobSecret.c, tt.msg)

		t.Run(testName, func(t *testing.T) {
			alice := NewAlice(tt.bob.BobPub)
			encoded := alice.Encode(tt.msg)
			log.Printf("ELGAMAL: encoded=%v", encoded.Value)
			decoded := tt.bob.Decode(encoded)
			log.Printf("ELGAMAL: decoded: %v", decoded)
			if !reflect.DeepEqual(decoded, tt.msg) {
				t.Errorf("decoded = %v, want %v", decoded, tt.msg)
			}
		})
	}
}
