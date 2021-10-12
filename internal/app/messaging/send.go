package messaging

import (
	"errors"
	"fmt"
	"log"
	"math/big"
	. "net"
)

type SendArgs struct {
	Addr string
	Msg  *big.Int
}

type WriteMsg = func(Msg, Conn) error

func SendMsg(addr string, msg Msg, write WriteMsg) error {
	conn, err := Dial("tcp", addr)
	if err != nil {
		errMsg := fmt.Sprintf("can't connect to %s: %v", addr, err)
		return errors.New(errMsg)
	}
	defer func() { _ = conn.Close() }()

	log.Printf("connected to %s", addr)
	err = write(msg, conn)
	if err != nil {
		return err
	}
	return nil
}
