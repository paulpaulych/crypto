package protocols

import (
	"errors"
	"fmt"
	"github.com/paulpaulych/crypto/shamir"
	"github.com/paulpaulych/crypto/shamir/cli/internal/messaging"
	"github.com/paulpaulych/crypto/shamir/cli/internal/tcp"
	"log"
	. "math/big"
	. "net"
)

func ShamirReader() messaging.ReadMsg {
	return read
}

func ShamirWriter(p *Int) messaging.WriteMsg {
	return func(msg messaging.Msg, conn Conn) error {
		return write(p, msg, conn)
	}
}

func write(P *Int, msg messaging.Msg, conn Conn) error {
	alice, err := shamir.InitAlice(P)
	if err != nil {
		log.Printf("failed to init alice: %d", err)
		return err
	}

	err = tcp.WriteBigIntWithLen(conn, P)
	if err != nil {
		errMsg := fmt.Sprintf("writing P failed: %v", err)
		return errors.New(errMsg)
	}

	step1out, err := alice.Step1(msg)
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

func read(conn Conn) (*Int, error) {
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
