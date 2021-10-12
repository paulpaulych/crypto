package protocols

import (
	"errors"
	"fmt"
	"github.com/paulpaulych/crypto/internal/app/messaging/msg-core"
	"github.com/paulpaulych/crypto/internal/app/messaging/nio"
	"math/big"
	. "net"
)

const (
	Shamir msg_core.ProtocolCode = iota
)

func ShamirWriter(p *big.Int) msg_core.Alice {
	return &writer{
		code:  Shamir,
		write: ShamirWriteFn(p),
	}
}

func ChooseBob(
	code msg_core.ProtocolCode,
	out func(Addr) nio.BlockWriter,
	onErr func(string),
) (msg_core.Bob, error) {
	switch code {
	case Shamir:
		return ShamirBob(out, onErr), nil
	default:
		msg := fmt.Sprintf("unknown protocol code %v", code)
		return nil, errors.New(msg)
	}
}

type writer struct {
	code  msg_core.ProtocolCode
	write func(nio.ByteReader, Conn) error
}

func (w writer) ProtocolCode() msg_core.ProtocolCode {
	return w.code
}
func (w writer) Write(msg nio.ByteReader, conn Conn) error {
	return w.write(msg, conn)
}
