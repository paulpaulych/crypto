package msg_core

import (
	"fmt"
	"io"
	"github.com/paulpaulych/crypto/internal/app/lang/nio"
)

type ProtocolCode = uint32


type Receiver interface {
	ProtocolCode() ProtocolCode
	RecvReader(src io.ReadWriter) (io.Reader, error)
}

type Sender interface {
	ProtocolCode() ProtocolCode
	SendWriter(target io.ReadWriter) (io.Writer, error)
}

type SendFunc func(target io.ReadWriter) (io.Writer, error)

type ReceiveFunc func(src io.ReadWriter) (io.Reader, error)

func EncryptedMsgWriter(target io.ReadWriter, s Sender) (io.Writer, error) {
	protocolCode := s.ProtocolCode()
	err := nio.WriteUint32(target, protocolCode)
	if err != nil {
		return nil, fmt.Errorf("failed to write protocol code %v: %v", protocolCode, err)
	}
	return s.SendWriter(target)
}

func EncryptedMsgReader(src io.ReadWriter, r Receiver) (io.Reader, error) {
	protocolCode, err := nio.ReadUint32(src)
	if err != nil {
		return nil, fmt.Errorf("failed to read protocol protocolCode: %s", err)
	}
	if protocolCode != r.ProtocolCode() {
		return nil, fmt.Errorf("unsupported protocol")
	}
	return r.RecvReader(src)
}


func NewReceiver(protocol ProtocolCode, fn ReceiveFunc) Receiver {
	return &receiver{protocol: protocol, fn: fn}
}
type receiver struct {
	protocol ProtocolCode
	fn ReceiveFunc
}
func (r receiver) ProtocolCode() ProtocolCode {
	return r.protocol
}
func (r receiver) RecvReader(src io.ReadWriter) (io.Reader, error) {
	return r.fn(src)
}

func NewSender(protocol ProtocolCode, fn SendFunc) Sender {
	return &sender{protocol: protocol, fn: fn}
}
type sender struct {
	protocol ProtocolCode
	fn SendFunc}
func (s sender) ProtocolCode() ProtocolCode {
	return s.protocol
}
func (s sender) SendWriter(src io.ReadWriter) (io.Writer, error) {
	return s.fn(src)
}