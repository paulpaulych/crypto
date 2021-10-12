package msg_core

import (
	"github.com/paulpaulych/crypto/internal/app/messaging/nio"
	"net"
)

type Msg = []byte

type ProtocolCode = uint32

type Bob = func(net.Conn)

type Alice interface {
	ProtocolCode() ProtocolCode
	Write(msg nio.ByteReader, conn net.Conn) error
}
