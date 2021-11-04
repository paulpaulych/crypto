package rsa

import (
	"bytes"
	"fmt"
	"io"
	"math/big"
	"log"
	"github.com/paulpaulych/crypto/internal/app/lang/nio"
	"github.com/paulpaulych/crypto/internal/app/messaging/msg-core"
	"github.com/paulpaulych/crypto/internal/core/rand"
	rsa "github.com/paulpaulych/crypto/internal/core/rsa-cipher"
)

// TODO: increase block size
const blockSize = 1

func SendFunc(bobPub []byte) (msg_core.SendFunc, error) {
	bp, err := readBobKey(bobPub)
	if err != nil {
		return nil, fmt.Errorf("can'r read bob public key: %v", err)
	}
	sendFunc := func(rw io.ReadWriter) (io.Writer, error) {
		target := &nio.BlockTarget {
			MetaWriter: rw,
			DataWriter: nio.WriterFunc(encrypt(bp, rw)),
		}
		writer, err := nio.NewBlockTransfer(blockSize).Writer(target)
		if err != nil {
			return nil, fmt.Errorf("error sending block: %v", err)
		}
		return writer, nil
	}

	return sendFunc, nil 
}

func ReceiveFunc(p, q *big.Int, onBobKey func([]byte) error) (msg_core.ReceiveFunc, error) {
	bob, err := rsa.NewBob(p, q, rand.CryptoSafeRandom())
	if err != nil {
		return nil, fmt.Errorf("bob init failed: %v", err)
	}
	bobPubBytes, err := writeBobKey(bob.BobPub)
	if err != nil {
		return nil, err
	}
	err = onBobKey(bobPubBytes)
	if err != nil {
		return nil, err
	}
	receiveFunc := func(rw io.ReadWriter) (io.Reader, error) {
		src := &nio.BlockSrc{
			MetaReader: rw,
			DataReader: nio.ReaderFunc(decrypt(bob, rw)),
		}
		blocks := nio.NewBlockTransfer(blockSize)
		return blocks.Reader(src), nil
	}
	return receiveFunc, nil
}

func encrypt(bp *rsa.BobPub, rw io.ReadWriter) nio.WriterFunc {
	max := bp.MaxValueCanBeEcnrypted()
	return func(p []byte) (int, error) {
		alice := rsa.NewAlice(bp)
		log.Println(fmtAlice(alice))

		given := new(big.Int).SetBytes(p)
		if given.Cmp(max) > 0 {
			return 0, fmt.Errorf("RSA: can't encrypt value %v beacause it's greater than %v", given, max)
		}
		encoded := alice.Encode(given)
		e := nio.WriteBigIntWithLen(rw, encoded.Value)
		if e != nil {
			return 0, fmt.Errorf("writing encoded msg failed: %v", e)
		}

		return len(p), nil
	}
}

func decrypt(b *rsa.Bob, rw io.ReadWriter) nio.ReaderFunc {
	return func(p []byte) (int, error) {
		encodedValue, err := nio.ReadBigIntWithLen(rw)
		if err == io.EOF {
			return 0, io.EOF
		}
		if err != nil {
			return 0, fmt.Errorf("can't read encoded: %v", err)
		}
		encoded := &rsa.Encoded{Value: encodedValue}
		b.Decode(encoded).FillBytes(p)
		return len(p), nil
	}
}

func fmtAlice(a *rsa.Alice) string {
	return fmt.Sprintln("Shamir node(Alice) initialized.\n",
		fmt.Sprintf("Bob public key: N=%v,d=%v", a.BobPub.N, a.BobPub.D),
	)
}

func writeBobKey(key *rsa.BobPub) ([]byte, error) {
	buf := &bytes.Buffer{}
	err := nio.WriteBigIntWithLen(buf, key.N)
	if err != nil {
		return nil, err
	}
	err = nio.WriteBigIntWithLen(buf, key.D)
	if err != nil {
		return nil, err
	}
	bytes := buf.Bytes()
	return bytes, nil
}

func readBobKey(b []byte) (*rsa.BobPub, error) {
	log.Println("rsa bob key bytes: ", b)
	buf := bytes.NewBuffer(b)
	N, err := nio.ReadBigIntWithLen(buf)
	if err != nil {
		return nil, err
	}
	D, err := nio.ReadBigIntWithLen(buf)
	if err != nil {
		return nil, err
	}
	return &rsa.BobPub{D: D, N: N}, nil
}
