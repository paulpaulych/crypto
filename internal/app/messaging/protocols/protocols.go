package protocols

import (
	dh "github.com/paulpaulych/crypto/internal/app/algorithms/diffie-hellman"
	"github.com/paulpaulych/crypto/internal/app/messaging/msg-core"
	"github.com/paulpaulych/crypto/internal/app/messaging/protocols/elgamal"
	"github.com/paulpaulych/crypto/internal/app/messaging/protocols/shamir"
	"github.com/paulpaulych/crypto/internal/app/nio"
	"io/ioutil"
	"math/big"
	. "net"
)

const (
	Shamir msg_core.ProtocolCode = iota
	Elgamal
)

func ShamirWriter(p *big.Int) msg_core.ConnWriter {
	return msg_core.NewConnWriter(Shamir, shamir.WriteFn(p))
}

func ShamirReader(out func(addr Addr) nio.ClosableWriter) msg_core.ConnReader {
	return msg_core.NewConnReader(Shamir, shamir.ReadFn(out))
}

func ElgamalWriter(cPub *dh.CommonPublicKey, bobPubFileName string) msg_core.ConnWriter {
	bytes, err := ioutil.ReadFile(bobPubFileName)
	if err != nil {
		panic(err)
	}
	bobPub := new(big.Int).SetBytes(bytes)
	return msg_core.NewConnWriter(Elgamal, elgamal.WriteFn(cPub, bobPub))
}

func ElgamalReader(cPub *dh.CommonPublicKey, out func(addr Addr) nio.ClosableWriter) msg_core.ConnReader {
	return msg_core.NewConnReader(Elgamal, elgamal.ReadFn(cPub, out))
}
