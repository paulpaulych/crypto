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

type SendArgs struct {
	Addr string
	Msg  *big.Int
	P    *big.Int
}

func SendShamirEncrypted(args *SendArgs) error {
	alice, err := shamir.InitAlice(args.P)
	if err != nil {
		log.Printf("failed to init alice: %d", err)
		return err
	}

	conn, err := net.Dial("tcp", args.Addr)
	if err != nil {
		errMsg := fmt.Sprintf("can't connect to %s: %v", args.Addr, err)
		return errors.New(errMsg)
	}
	defer func() { _ = conn.Close() }()

	log.Printf("connected to %s", args.Addr)

	err = tcp.WriteBigIntWithLen(conn, args.P)
	if err != nil {
		errMsg := fmt.Sprintf("writing P failed: %v", err)
		return errors.New(errMsg)
	}

	step1out, err := alice.Step1(args.Msg)
	if err != nil {
		errMsg := fmt.Sprintf("writing step1out failed: %v", err)
		return errors.New(errMsg)
	}

	err = tcp.WriteBigIntWithLen(conn, step1out)
	if err != nil {
		errMsg := fmt.Sprintf("writing step1out failed: %v", err)
		return errors.New(errMsg)
	}

	step2out, err := tcp.ReadBigIntWithLen(conn)
	if err != nil {
		errMsg := fmt.Sprintf("reading step2out failed: %v", err)
		return errors.New(errMsg)
	}

	step3out := alice.Step3(step2out)
	err = tcp.WriteBigIntWithLen(conn, step3out)
	if err != nil {
		errMsg := fmt.Sprintf("writing step3out failed: %v", err)
		return errors.New(errMsg)
	}

	return nil
}
