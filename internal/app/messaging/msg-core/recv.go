package msg_core

import (
	"github.com/paulpaulych/crypto/internal/app/nio"
	"log"
	. "net"
)

func RecvMessage(reader ConnReader) func(Conn) {
	return func(conn Conn) {
		defer func() { _ = conn.Close() }()

		protocolCode, err := nio.ReadUint32(conn)
		if err != nil {
			log.Printf("failed to read protocol protocolCode: %s", err)
			return
		}
		if protocolCode != reader.ProtocolCode() {
			log.Printf("unsupprted protocol")
			return
		}
		err = reader.MsgReader()(conn)
		if err != nil {
			log.Printf("error reading message: %s", err)
			return
		}
	}
}
