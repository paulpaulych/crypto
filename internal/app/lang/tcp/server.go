package tcp

import (
	"fmt"
	"log"
	"net"
)

func StartServer(
	bindAddr string,
	errs chan<- error,
	handler func(net.Conn) error,
) error {
	listener, err := net.Listen("tcp", bindAddr)
	if err != nil {
		return fmt.Errorf("can't bind to %s: %v", bindAddr, err)
	}
	defer func() { _ = listener.Close() }()

	for {
		log.Printf("TCP_SERVER: waiting for connection on %v", bindAddr)
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("TCP_SERVER: can't accept connection: %v", err)
		}
		go func() {
			err := handler(conn)
			if err != nil {
				errs <- err
			}
		}()
	}
}
