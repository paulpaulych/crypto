package msg_core

import (
	"errors"
	"fmt"
	"github.com/paulpaulych/crypto/internal/app/tcp"
	"log"
	"math/big"
	. "net"
)

type SendArgs struct {
	Addr string
	Msg  *big.Int
}

func SendMsg(addr string, msg Msg, writer MsgWriter) error {
	conn, err := Dial("tcp", addr)
	if err != nil {
		errMsg := fmt.Sprintf("can't connect to %s: %v", addr, err)
		return errors.New(errMsg)
	}
	defer func() { _ = conn.Close() }()

	protocolCode := writer.ProtocolCode()
	err = tcp.WriteUint32(conn, protocolCode)
	if err != nil {
		errMsg := fmt.Sprintf("failed to write protocol protocolCode %v: %v", protocolCode, err)
		return errors.New(errMsg)
	}

	log.Printf("connected to %s", addr)
	err = writer.Write(msg, conn)
	if err != nil {
		return err
	}
	return nil
}
