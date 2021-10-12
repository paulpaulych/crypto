package tcp

import (
	"errors"
	"fmt"
	"log"
	"net"
)

func StartServer(bindAddr string, handler func(net.Conn)) error {
	listener, err := net.Listen("tcp", bindAddr)
	if err != nil {
		errMsg := fmt.Sprintf("can't bind to %s: %v", bindAddr, err)
		return errors.New(errMsg)
	}
	defer func() { _ = listener.Close() }()

	for {
		log.Printf("waiting for connection on %v", bindAddr)
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("can't accept connection: %v", err)
		}
		handler(conn)
		err = conn.Close()
		if err != nil {
			log.Printf("failed to close connection: %v", err)
		}
	}
}
