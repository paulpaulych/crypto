package shamir

import (
	"fmt"
	"github.com/paulpaulych/crypto/internal/app/algorithms/shamir-cipher"
	"github.com/paulpaulych/crypto/internal/app/nio"
	"io"
	"log"
	. "math/big"
	. "net"
)

const blockSize = 4

func WriteFn(p *Int) func(msg io.Reader, conn Conn) error {
	return func(msg io.Reader, conn Conn) error {
		err := nio.WriteBigIntWithLen(conn, p)
		if err != nil {
			return fmt.Errorf("writing P failed: %v", err)
		}

		err = nio.NewBlockTransfer(blockSize).WriteBlocks(nio.WriteProps{
			From:       msg,
			MetaWriter: conn,
			DataWriter: nio.NewFnWriter(encoder(p, conn)),
		})
		if err != nil {
			return fmt.Errorf("error sending block: %v", err)
		}
		return nil
	}
}

func encoder(p *Int, conn Conn) func([]byte) error {
	return func(block []byte) error {
		alice, err := shamir_cipher.InitAlice(p)
		if err != nil {
			log.Printf("failed to init alice: %d", err)
			return err
		}

		msgInt := new(Int).SetBytes(block)
		step1out, err := alice.Step1(msgInt)
		if err != nil {
			return fmt.Errorf("writing step1out failed: %v", err)
		}

		err = nio.WriteBigIntWithLen(conn, step1out)
		if err != nil {
			return fmt.Errorf("writing step1out failed: %v", err)
		}

		step2out, err := nio.ReadBigIntWithLen(conn)
		if err != nil {
			return fmt.Errorf("reading step2out failed: %v", err)
		}

		step3out := alice.Step3(step2out)
		err = nio.WriteBigIntWithLen(conn, step3out)
		if err != nil {
			return fmt.Errorf("writing step3out failed: %v", err)
		}

		return nil
	}
}

func ReadFn(
	output func(Addr) nio.ClosableWriter,
) func(conn Conn) error {
	return func(conn Conn) error {
		out := output(conn.RemoteAddr())
		defer func() {
			err := out.Close()
			if err != nil {
				log.Printf("failed to close writer: %s", err)
			}
		}()

		p, err := nio.ReadBigIntWithLen(conn)
		if err != nil {
			return fmt.Errorf("can't read p: %s", err)
		}

		err = nio.NewBlockTransfer(blockSize).
			ReadBlocks(nio.ReadProps{
				MetaReader: conn,
				DataReader: nio.NewFnReader(decoder(p, conn)),
				To:         out,
			})
		if err != nil {
			return fmt.Errorf("can't transfer: %s", err)
		}
		return nil
	}
}

func decoder(p *Int, conn Conn) func(buf []byte) (int, error) {
	return func(buf []byte) (int, error) {
		bob, err := shamir_cipher.InitBob(p)
		if err != nil {
			return 0, fmt.Errorf("failed to init bob: %d", err)
		}

		step1out, err := nio.ReadBigIntWithLen(conn)
		if err == io.EOF {
			return 0, io.EOF
		}
		if err != nil {
			return 0, fmt.Errorf("can't read step1out: %v", err)
		}

		step2out := bob.Step2(step1out)
		err = nio.WriteBigIntWithLen(conn, step2out)
		if err != nil {
			return 0, fmt.Errorf("can't write step2out: %v", err)
		}

		step3out, err := nio.ReadBigIntWithLen(conn)
		if err != nil {
			return 0, fmt.Errorf("can't write step2out: %v", err)
		}
		bob.Decode(step3out).FillBytes(buf)
		return blockSize, nil
	}
}
