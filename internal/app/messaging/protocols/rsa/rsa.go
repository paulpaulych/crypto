package rsa

import (
	"fmt"
	"github.com/paulpaulych/crypto/internal/app/algorithms/rand"
	"github.com/paulpaulych/crypto/internal/app/algorithms/rsa-cipher"
	"github.com/paulpaulych/crypto/internal/app/nio"
	"io"
	"log"
	. "math/big"
	. "net"
	"os"
)

const bobKeyFileName = "bob_rsa.key"

// TODO: increase block size
const blockSize = 1

func WriteFn(bobPubFileName string) func(msg io.Reader, conn Conn) error {
	return func(msg io.Reader, conn Conn) error {
		file, err := os.Open(bobPubFileName)
		if err != nil {
			return err
		}
		bobPub, err := readBobKey(file)
		if err != nil {
			return err
		}

		err = nio.NewBlockTransfer(blockSize).WriteBlocks(nio.WriteProps{
			From:       msg,
			MetaWriter: conn,
			DataWriter: nio.NewFnWriter(encoder(bobPub, conn)),
		})
		if err != nil {
			return fmt.Errorf("error sending block: %v", err)
		}
		return nil
	}
}

func encoder(bobPub *rsa_cipher.BobPub, conn Conn) func([]byte) error {
	return func(block []byte) error {
		alice := rsa_cipher.NewAlice(bobPub)
		log.Println(fmtAlice(alice))

		msgInt := new(Int).SetBytes(block)
		encoded := alice.Encode(msgInt)

		err := nio.WriteBigIntWithLen(conn, encoded.Value)
		if err != nil {
			return fmt.Errorf("writing encoded msg failed: %v", err)
		}

		return nil
	}
}

func ReadFn(
	p *Int,
	g *Int,
	output func(Addr) nio.ClosableWriter,
) func(conn Conn) error {
	bob, err := rsa_cipher.NewBob(p, g, rand.CryptoSafeRandom())
	if err != nil {
		log.Panicf("bob init failed: %v", err)
	}
	fmt.Println(fmtBob(bob))
	return func(conn Conn) error {
		out := output(conn.RemoteAddr())
		defer func() {
			err := out.Close()
			if err != nil {
				log.Printf("failed to close writer: %s", err)
			}
		}()

		err := nio.NewBlockTransfer(blockSize).
			ReadBlocks(nio.ReadProps{
				MetaReader: conn,
				DataReader: nio.NewFnReader(decoder(bob, conn)),
				To:         out,
			})
		if err != nil {
			return fmt.Errorf("can't transfer: %s", err)
		}
		return nil
	}
}

func decoder(bob *rsa_cipher.Bob, conn Conn) func(buf []byte) (int, error) {
	return func(buf []byte) (int, error) {
		encodedValue, err := nio.ReadBigIntWithLen(conn)
		if err == io.EOF {
			return 0, io.EOF
		}
		if err != nil {
			return 0, fmt.Errorf("can't read encoded: %v", err)
		}
		encoded := &rsa_cipher.Encoded{Value: encodedValue}
		bob.Decode(encoded).FillBytes(buf)
		return blockSize, nil
	}
}

func fmtAlice(a *rsa_cipher.Alice) string {
	return fmt.Sprintln("Shamir node(Alice) initialized.\n",
		fmt.Sprintf("Bob public key: N=%v,d=%v", a.BobPub.N, a.BobPub.D),
	)
}

func fmtBob(bob *rsa_cipher.Bob) string {
	file := nio.NewFileWriter(bobKeyFileName, func() {
		fmt.Println("PUBLIC KEY SAVED TO", bobKeyFileName)
	})
	defer func() { _ = file.Close() }()
	err := writeBobKey(bob.BobPub, file)
	if err != nil {
		log.Panicf("failed to save key to file: %v", err)
	}
	return fmt.Sprintln("Shamir node(Bob) initialized.\n",
		fmt.Sprintf("Public key: N=%v, d=%v\n", bob.BobPub.N, bob.BobPub.D),
	)
}

func writeBobKey(key *rsa_cipher.BobPub, writer io.Writer) error {
	err := nio.WriteBigIntWithLen(writer, key.N)
	if err != nil {
		return err
	}
	err = nio.WriteBigIntWithLen(writer, key.D)
	if err != nil {
		return err
	}
	return nil
}

func readBobKey(reader io.Reader) (*rsa_cipher.BobPub, error) {
	N, err := nio.ReadBigIntWithLen(reader)
	if err != nil {
		return nil, err
	}
	D, err := nio.ReadBigIntWithLen(reader)
	if err != nil {
		return nil, err
	}
	return &rsa_cipher.BobPub{D: D, N: N}, nil
}
