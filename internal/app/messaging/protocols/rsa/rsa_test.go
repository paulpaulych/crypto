package rsa

import (
	"fmt"
	"io"
	"math/big"
	"reflect"
	"testing"

	"github.com/paulpaulych/crypto/internal/app/lang/nio"
	"github.com/paulpaulych/crypto/internal/app/lang/rw-test"
	"github.com/paulpaulych/crypto/internal/core/rand"
	rsa "github.com/paulpaulych/crypto/internal/core/rsa-cipher"
)

func Test_protocol(t *testing.T) {
	bob, e := rsa.NewBob(big.NewInt(30803), big.NewInt(1297), rand.ConstRand(big.NewInt(17)))
	if e != nil {
		t.Errorf("test init error: %v", e)
	}
	type test struct {
		name       string
		given      []byte
		wantRead   []byte
		wantReadE  error
		wantWriteE error
	}
	tests := []test{
		{
			name:     "ok - single byte",
			given:    []byte{2},
			wantRead: []byte{2},
		},
		{
			name:     "ok - multiple bytes",
			given:    []byte{0, 0, 0, 0, 0, 0, 1},
			wantRead: []byte{0, 0, 0, 0, 0, 0, 1},
		},
		{
			name:     "ok - multiple bytes",
			given:    []byte{1, 0, 0, 0},
			wantRead: []byte{1, 0, 0, 0},
		},
		{
			name:     "ok - empty array",
			given:    make([]byte, 0),
			wantRead: make([]byte, 0),
		},
	}
	for _, tt := range tests {
		type readResult struct {
			bytes []byte
			n     int
		}
		sender := wr_test.Peer{
			Act: func(rw io.ReadWriter) (wr_test.Result, error) {
				writer := nio.WriterFunc(encrypt(bob.BobPub, rw))
				return writer.Write(tt.given)
			},
			Check: func(got wr_test.Result, gotE error) error {
				if gotE != tt.wantWriteE {
					return fmt.Errorf("encrypt().Write() got error = %v, want = %v", gotE, tt.wantWriteE)
				}
				if got != len(tt.given) {
					return fmt.Errorf("encrypt().Write() = %v, want %v", got, len(tt.given))
				}
				return nil
			},
		}
		receiver := wr_test.Peer{
			Act: func(rw io.ReadWriter) (wr_test.Result, error) {
				reader := nio.ReaderFunc(decrypt(bob, rw))
				buf := make([]byte, len(tt.wantRead))
				n, err := reader.Read(buf)
				return readResult{bytes: buf, n: n}, err
			},
			Check: func(got wr_test.Result, gotE error) error {
				if gotE != tt.wantReadE {
					return fmt.Errorf("decrypt().Read() got error = %v, want = %v", gotE, tt.wantReadE)
				}
				g, _ := got.(readResult)
				if g.n != len(tt.wantRead) {
					return fmt.Errorf("expected read %v bytes, actual %v", len(tt.wantRead), g.n)
				}
				if !reflect.DeepEqual(g.bytes, tt.wantRead) {
					return fmt.Errorf("read() = %v, want %v", g.bytes, tt.wantRead)
				}
				return nil
			},
		}
		wr_test.RunReadWriteTest(t, &wr_test.ReadWriteTest{
			Name: tt.name,
			A:    sender,
			B:    receiver,
		})
	}
}
