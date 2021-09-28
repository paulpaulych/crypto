package main

import (
	"errors"
	"fmt"
	"net"
)

func tcpSend(addr string, bytes []byte) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		errMsg := fmt.Sprintf("can't connect to %s: %v", addr, err)
		return errors.New(errMsg)
	}
	_, err = conn.Write(bytes)
	if err != nil {
		errMsg := fmt.Sprintf("error sending to %s: %v", addr, err)
		return errors.New(errMsg)
	}
	return nil
}
