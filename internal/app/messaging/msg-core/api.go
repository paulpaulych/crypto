package msg_core

import (
	"io"
	. "net"
)

type Msg = []byte
type ProtocolCode = uint32

type ConnReader interface {
	ProtocolCode() uint32
	MsgReader() func(Conn) error
}

func NewConnReader(
	protocol ProtocolCode,
	msgReader func(Conn) error,
) ConnReader {
	return &reader{protocol: protocol, msgReader: msgReader}
}

type ConnWriter interface {
	ProtocolCode() ProtocolCode
	Write(msg io.Reader, conn Conn) error
}

type ConnWriteFn = func(msg io.Reader, conn Conn) error

func NewConnWriter(
	protocol ProtocolCode,
	msgWriter ConnWriteFn,
) ConnWriter {
	return &writer{protocol: protocol, msgWriter: msgWriter}
}

type reader struct {
	protocol  ProtocolCode
	msgReader func(Conn) error
}

func (r reader) ProtocolCode() ProtocolCode {
	return r.protocol
}
func (r reader) MsgReader() func(conn Conn) error {
	return r.msgReader
}

type writer struct {
	protocol  ProtocolCode
	msgWriter func(io.Reader, Conn) error
}

func (w writer) ProtocolCode() ProtocolCode {
	return w.protocol
}
func (w writer) Write(msg io.Reader, conn Conn) error {
	return w.msgWriter(msg, conn)
}
