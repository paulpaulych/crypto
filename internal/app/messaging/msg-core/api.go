package msg_core

import (
	"io"
	"net"
)

type Msg = []byte

type ProtocolCode = uint32

type Bob = func(net.Conn)

type Alice interface {
	ProtocolCode() ProtocolCode
	Write(msg io.Reader, conn net.Conn) error
}
