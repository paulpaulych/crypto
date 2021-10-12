package msg_core

import (
	"math/big"
	"net"
)

type Msg = *big.Int

type ProtocolCode = uint32

type Read = func(net.Conn) (Msg, error)

type MsgWriter interface {
	ProtocolCode() ProtocolCode
	Write(Msg, net.Conn) error
}
