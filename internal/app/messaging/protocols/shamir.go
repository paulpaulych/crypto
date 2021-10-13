package protocols

import (
	"errors"
	"fmt"
	"github.com/paulpaulych/crypto/internal/app/algorithms/shamir"
	"github.com/paulpaulych/crypto/internal/app/messaging/msg-core"
	"github.com/paulpaulych/crypto/internal/app/messaging/nio"
	"github.com/paulpaulych/crypto/internal/app/tcp"
	"io"
	"log"
	. "math/big"
	. "net"
)

const blockSize = 4

func ShamirWriteFn(p *Int) func(msg io.Reader, conn Conn) error {
	return func(msg io.Reader, conn Conn) error {
		err := tcp.WriteBigIntWithLen(conn, p)
		if err != nil {
			errMsg := fmt.Sprintf("writing P failed: %v", err)
			return errors.New(errMsg)
		}
		for {
			buf := make([]byte, blockSize)
			_, err := msg.Read(buf)
			if err == io.EOF {
				return nil
			}
			if err != nil {
				errMsg := fmt.Sprintf("error reading message: %v", err)
				return errors.New(errMsg)
			}
			err = sendBlock(p, conn, buf)
			if err != nil {
				errMsg := fmt.Sprintf("error sending block: %v", err)
				return errors.New(errMsg)
			}
		}
	}
}

func sendBlock(p *Int, conn Conn, block []byte) error {
	alice, err := shamir.InitAlice(p)
	if err != nil {
		log.Printf("failed to init alice: %d", err)
		return err
	}

	msgInt := new(Int).SetBytes(block)
	step1out, err := alice.Step1(msgInt)
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

func ShamirBob(
	output func(Addr) nio.ClosableWriter,
	onErr func(string),
) msg_core.Bob {
	return func(conn Conn) {
		out := output(conn.RemoteAddr())
		defer func() {
			err := out.Close()
			if err != nil {
				log.Printf("failed to close writer: %s", err)
			}
		}()

		p, err := tcp.ReadBigIntWithLen(conn)
		if err != nil {
			onErr(fmt.Sprintf("can't read p: %v", err))
			return
		}

		blockReader := &blockReader{p: p, conn: conn}
		blockBuf := make([]byte, blockSize)
		for {
			block, err := blockReader.Read(blockBuf)
			if err == io.EOF {
				return
			}
			if err != nil {
				onErr(fmt.Sprintf("can't read block: %v", err))
				return
			}
			log.Printf("received block: len=%v, bytes=%v", len(blockBuf), block)
			_, err = out.Write(blockBuf)
			if err != nil {
				onErr(fmt.Sprintf("error writing received message: %v", err))
				return
			}
		}
	}
}

type blockReader struct {
	p    *Int
	conn Conn
}

func (r blockReader) Read(buf []byte) (uint, error) {
	bob, err := shamir.InitBob(r.p)
	if err != nil {
		return 0, errors.New(fmt.Sprintf("failed to init bob: %d", err))
	}

	step1out, err := tcp.ReadBigIntWithLen(r.conn)
	if err == io.EOF {
		return 0, io.EOF
	}
	if err != nil {
		return 0, errors.New(fmt.Sprintf("can't read step1out: %v", err))
	}

	step2out := bob.Step2(step1out)
	err = tcp.WriteBigIntWithLen(r.conn, step2out)
	if err != nil {
		return 0, errors.New(fmt.Sprintf("can't write step2out: %v", err))
	}

	step3out, err := tcp.ReadBigIntWithLen(r.conn)
	if err != nil {
		return 0, errors.New(fmt.Sprintf("can't write step2out: %v", err))
	}
	bob.Decode(step3out).FillBytes(buf)
	return blockSize, nil
}
