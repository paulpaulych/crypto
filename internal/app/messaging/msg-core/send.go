package msg_core

import (
	"fmt"
	"github.com/paulpaulych/crypto/internal/app/nio"
	"io"
	"log"
	"math/big"
	. "net"
)

type SendArgs struct {
	Addr string
	Msg  *big.Int
}

func SendMsg(addr string, msg io.Reader, alice ConnWriter) error {
	conn, err := Dial("tcp", addr)
	if err != nil {
		return fmt.Errorf("can't connect to %s: %v", addr, err)
	}
	defer func() { _ = conn.Close() }()

	protocolCode := alice.ProtocolCode()
	err = nio.WriteUint32(conn, protocolCode)
	if err != nil {
		return fmt.Errorf("failed to write protocol code %v: %v", protocolCode, err)
	}

	log.Printf("connected to %s", addr)
	err = alice.Write(msg, conn)
	if err != nil {
		return err
	}
	return nil
}
