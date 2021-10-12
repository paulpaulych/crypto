package protocols

import (
	"errors"
	"fmt"
	"github.com/paulpaulych/crypto/internal/app/messaging/msg-core"
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
	out func(Addr) MsgWriter,
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
	write func(msg_core.Msg, Conn) error
}

func (w writer) ProtocolCode() msg_core.ProtocolCode {
	return w.code
}
func (w writer) Write(msg msg_core.Msg, conn Conn) error {
	return w.write(msg, conn)
}
