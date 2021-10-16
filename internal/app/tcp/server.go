package tcp

import (
	"fmt"
	"log"
	"net"
)

func StartServer(bindAddr string, handler func(net.Conn)) error {
	listener, err := net.Listen("tcp", bindAddr)
	if err != nil {
		return fmt.Errorf("can't bind to %s: %v", bindAddr, err)
	}
	defer func() { _ = listener.Close() }()

	for {
		log.Printf("waiting for connection on %v", bindAddr)
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("can't accept connection: %v", err)
		}
		go handler(conn)
	}
}
