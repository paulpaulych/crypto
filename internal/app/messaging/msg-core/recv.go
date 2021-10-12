package msg_core

import (
	"errors"
	"fmt"
	"github.com/paulpaulych/crypto/internal/app/tcp"
	"log"
	. "math/big"
	. "net"
)

type GetReader = func(code ProtocolCode) (Read, error)

func ListenForMsg(bindAddr string, reader GetReader) error {
	listener, err := Listen("tcp", bindAddr)
	if err != nil {
		errMsg := fmt.Sprintf("can't bind to %s: %v", bindAddr, err)
		return errors.New(errMsg)
	}
	defer func() { _ = listener.Close() }()

	for {
		log.Printf("waiting for connection on %v", bindAddr)
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("can't accept connection: %v\n", err)
		}
		code, err := tcp.ReadUint32(conn)
		if err != nil {
			log.Printf("failed to read protocol code: %s. Continuing listening", err)
			continue
		}
		read, err := reader(code)
		if err != nil {
			log.Printf("")
		}

		msg, err := read(conn)
		if err != nil {
			log.Printf("can't read message: %v\n", err)
			continue
		}
		printMsg(msg, conn.RemoteAddr().String())
	}
}

func printMsg(msg *Int, from string) {
	fmt.Printf("RECEIVED MESSAGE from %s: %s", from, msg)
	fmt.Println()
}
