package protocols

import (
	"fmt"
	"github.com/paulpaulych/crypto/internal/app/messaging/msg-core"
	"github.com/paulpaulych/crypto/internal/app/messaging/nio"
	"io"
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
	out func(Addr) nio.ClosableWriter,
	onErr func(error),
) (msg_core.Bob, error) {
	switch code {
	case Shamir:
		return ShamirBob(out, onErr), nil
	default:
		return nil, fmt.Errorf("unknown protocol code %v", code)
	}
}

type writer struct {
	code  msg_core.ProtocolCode
	write func(io.Reader, Conn) error
}

func (w writer) ProtocolCode() msg_core.ProtocolCode {
	return w.code
}
func (w writer) Write(msg io.Reader, conn Conn) error {
	return w.write(msg, conn)
}
