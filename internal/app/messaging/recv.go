package messaging

import (
	"errors"
	"fmt"
	"log"
	. "math/big"
	. "net"
)

type ReadMsg = func(Conn) (*Int, error)

func ListenForMsg(bindAddr string, readMsg ReadMsg) error {
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
		msg, err := readMsg(conn)
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
