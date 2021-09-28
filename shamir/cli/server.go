package main

import (
	"errors"
	"fmt"
	"log"
	"net"
)

func startServer(
	addr string,
	bufSize uint,
	onRecv func([]byte),
) error {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		errMsg := fmt.Sprintf("can't bind to %s: %v", addr, err)
		return errors.New(errMsg)
	}
	buf := make([]byte, bufSize)
	for {
		log.Printf("waiting for connection on %v", addr)
		conn, err := listener.Accept()
		if err != nil {
			errMsg := fmt.Sprintf("can't accept connection %s: %v", addr, err)
			return errors.New(errMsg)
		}
		size, err := conn.Read(buf)
		if err != nil {
			errMsg := fmt.Sprintf("can't read bytes: %v", err)
			return errors.New(errMsg)
		}
		onRecv(buf[:size])
	}
}
