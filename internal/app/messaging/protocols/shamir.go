package protocols

import (
	"errors"
	"fmt"
	"github.com/paulpaulych/crypto/internal/app/algorithms/shamir"
	"github.com/paulpaulych/crypto/internal/app/messaging/msg-core"
	"github.com/paulpaulych/crypto/internal/app/messaging/nio"
	"github.com/paulpaulych/crypto/internal/app/tcp"
	"log"
	. "math/big"
	. "net"
)

const blockSize = 4

func ShamirWriteFn(p *Int) func(msg nio.ByteReader, conn Conn) error {
	return func(msg nio.ByteReader, conn Conn) error {
		err := tcp.WriteBigIntWithLen(conn, p)
		if err != nil {
			errMsg := fmt.Sprintf("writing P failed: %v", err)
			return errors.New(errMsg)
		}
		totalBytes, err := msg.TotalBytes()
		if err != nil {
			errMsg := fmt.Sprintf("error counting total message size: %v", err)
			return errors.New(errMsg)
		}
		totalBlocks := divUp(totalBytes, blockSize)
		err = tcp.WriteUint32(conn, uint32(totalBlocks))
		if err != nil {
			errMsg := fmt.Sprintf("writing totalBlocks failed: %v", err)
			return errors.New(errMsg)
		}

		for {
			page, err := msg.Read(blockSize)
			log.Printf("sending block: len=%v, bytes=%v, hasMore=%v", len(page.Bytes), page.Bytes, page.HasMore)
			if err != nil {
				errMsg := fmt.Sprintf("error reading message: %v", err)
				return errors.New(errMsg)
			}
			err = sendBlock(p, conn, page.Bytes)
			if err != nil {
				errMsg := fmt.Sprintf("error sending block: %v", err)
				return errors.New(errMsg)
			}
			if !page.HasMore {
				return nil
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
	output func(Addr) nio.BlockWriter,
	onErr func(string),
) msg_core.Bob {
	return func(conn Conn) {
		out := output(conn.RemoteAddr())
		p, err := tcp.ReadBigIntWithLen(conn)
		if err != nil {
			onErr(fmt.Sprintf("can't read p: %v", err))
			return
		}
		totalBlocks, err := tcp.ReadUint32(conn)
		if err != nil {
			onErr(fmt.Sprintf("can't read totalBlocks: %v", err))
			return
		}
		log.Printf("totalBlocks: %v", totalBlocks)
		for i := uint32(0); i < totalBlocks; i++ {
			block, err := recvBlock(p, conn)
			if err != nil {
				onErr(fmt.Sprintf("can't read block: %v", err))
				return
			}
			hasMore := i != totalBlocks-1
			log.Printf("received block: len=%v, bytes=%v, hasMore=%v", len(block), block, hasMore)
			err = out.Write(block, hasMore)
			if err != nil {
				onErr(fmt.Sprintf("error writing received message: %v", err))
				return
			}
		}
	}
}

func recvBlock(p *Int, conn Conn) ([]byte, error) {
	bob, err := shamir.InitBob(p)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("failed to init bob: %d", err))
	}

	step1out, err := tcp.ReadBigIntWithLen(conn)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("can't read step1out: %v", err))
	}

	step2out := bob.Step2(step1out)
	err = tcp.WriteBigIntWithLen(conn, step2out)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("can't write step2out: %v", err))
	}

	step3out, err := tcp.ReadBigIntWithLen(conn)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("can't write step2out: %v", err))
	}

	bytes := bob.Decode(step3out).Bytes()
	return bytes, nil
}

func divUp(x uint, y uint) uint {
	return (x + y - 1) / y
}
