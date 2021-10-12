package protocols

import (
	"errors"
	"fmt"
	"github.com/paulpaulych/crypto/internal/app/algorithms/shamir"
	"github.com/paulpaulych/crypto/internal/app/messaging/msg-core"
	"github.com/paulpaulych/crypto/internal/app/tcp"
	"log"
	. "math/big"
	. "net"
)

func ShamirWriteFn(p *Int) func(msg msg_core.Msg, conn Conn) error {
	return func(msg msg_core.Msg, conn Conn) error {
		alice, err := shamir.InitAlice(p)
		if err != nil {
			log.Printf("failed to init alice: %d", err)
			return err
		}

		err = tcp.WriteBigIntWithLen(conn, p)
		if err != nil {
			errMsg := fmt.Sprintf("writing P failed: %v", err)
			return errors.New(errMsg)
		}

		step1out, err := alice.Step1(new(Int).SetBytes(msg))
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
}

type MsgWriter interface {
	Write(p []byte, hasMore bool) error
}

func ShamirBob(
	output func(Addr) MsgWriter,
	onErr func(string),
) msg_core.Bob {
	return func(conn Conn) {
		out := output(conn.RemoteAddr())
		p, err := tcp.ReadBigIntWithLen(conn)
		if err != nil {
			onErr(fmt.Sprintf("can't read p: %v", err))
			return
		}

		bob, err := shamir.InitBob(p)
		if err != nil {
			onErr(fmt.Sprintf("failed to init bob: %d", err))
			return
		}

		step1out, err := tcp.ReadBigIntWithLen(conn)
		if err != nil {
			onErr(fmt.Sprintf("can't read step1out: %v", err))
			return
		}

		step2out := bob.Step2(step1out)
		err = tcp.WriteBigIntWithLen(conn, step2out)
		if err != nil {
			onErr(fmt.Sprintf("can't write step2out: %v", err))
			return
		}

		step3out, err := tcp.ReadBigIntWithLen(conn)
		if err != nil {
			onErr(fmt.Sprintf("can't write step2out: %v", err))
			return
		}

		msgBytes := bob.Decode(step3out).Bytes()
		err = out.Write(msgBytes, false)
		if err != nil {
			onErr(fmt.Sprintf("error writing received message: %v", err))
			return
		}
	}
}
