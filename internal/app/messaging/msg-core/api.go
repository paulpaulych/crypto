package msg_core

import (
	"net"
)

type Msg = []byte

type ProtocolCode = uint32

type Bob = func(net.Conn)

type Alice interface {
	ProtocolCode() ProtocolCode
	Write(Msg, net.Conn) error
}
