package protocols

import (
	"errors"
	"fmt"
	"github.com/paulpaulych/crypto/internal/app/messaging/msg-core"
	"math/big"
	"net"
)

const (
	Shamir msg_core.ProtocolCode = iota
)

func ShamirWriter(p *big.Int) msg_core.MsgWriter {
	return &writer{
		code:  Shamir,
		write: ShamirWriteFn(p),
	}
}

func GetProtocolReader(code msg_core.ProtocolCode) (msg_core.Read, error) {
	switch code {
	case Shamir:
		return ShamirReader(), nil
	default:
		msg := fmt.Sprintf("unknown protocol code %v", code)
		return nil, errors.New(msg)
	}
}

type writer struct {
	code  msg_core.ProtocolCode
	write func(msg_core.Msg, net.Conn) error
}

func (w writer) ProtocolCode() msg_core.ProtocolCode {
	return w.code
}
func (w writer) Write(msg msg_core.Msg, conn net.Conn) error {
	return w.write(msg, conn)
}
