package command

import (
	"errors"
	"fmt"
	"github.com/paulpaulych/crypto/shamir"
	"github.com/paulpaulych/crypto/shamir/cli/tcp"
	"log"
	"math/big"
	"net"
)

type RecvArgs struct {
	BindAddr string
}

func RecvMessages(args *RecvArgs) error {
	listener, err := net.Listen("tcp", args.BindAddr)
	defer func() { _ = listener.Close() }()

	if err != nil {
		errMsg := fmt.Sprintf("can't bind to %s: %v", args.BindAddr, err)
		return errors.New(errMsg)
	}

	for {
		log.Printf("waiting for connection on %v", args.BindAddr)
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("can't accept connection: %v\n", err)
		}
		msg, err := recvEncryptedMsg(conn)
		if err != nil {
			log.Printf("can't receive message: %v\n", err)
			continue
		}
		printMsg(msg)
	}
}

func recvEncryptedMsg(conn net.Conn) (*big.Int, error) {
	p, err := tcp.ReadBigIntWithLen(conn)
	if err != nil {
		errMsg := fmt.Sprintf("can't read p: %v", err)
		return nil, errors.New(errMsg)
	}

	bob, err := shamir.InitBob(p)
	if err != nil {
		log.Printf("failed to init bob: %d", err)
		return nil, err
	}

	step1out, err := tcp.ReadBigIntWithLen(conn)
	if err != nil {
		errMsg := fmt.Sprintf("can't read step1out: %v", err)
		return nil, errors.New(errMsg)
	}

	step2out := bob.Step2(step1out)
	err = tcp.WriteBigIntWithLen(conn, step2out)
	if err != nil {
		errMsg := fmt.Sprintf("can't write step2out: %v", err)
		return nil, errors.New(errMsg)
	}

	step3out, err := tcp.ReadBigIntWithLen(conn)
	if err != nil {
		errMsg := fmt.Sprintf("can't read step3out: %v", err)
		return nil, errors.New(errMsg)
	}

	msg := bob.Decode(step3out)
	return msg, nil
}

func printMsg(msg *big.Int) {
	fmt.Printf("received message: %s", msg)
	fmt.Println()
}
