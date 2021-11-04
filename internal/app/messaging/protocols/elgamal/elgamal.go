package elgamal

import (
	"errors"
	"fmt"
	"github.com/paulpaulych/crypto/internal/app/lang/nio"
	"github.com/paulpaulych/crypto/internal/app/messaging/msg-core"
	dh "github.com/paulpaulych/crypto/internal/core/diffie-hellman"
	"github.com/paulpaulych/crypto/internal/core/elgamal-cipher"
	"github.com/paulpaulych/crypto/internal/core/rand"
	"io"
	"io/ioutil"
	"log"
	"math/big"
)

const bobPubKeyFile = "bob_elgamal.key"

// TODO increase block size
const blockSize = 1

func SendFunc(
	cp dh.CommonPublicKey,
	bobPub *big.Int,
) msg_core.SendFunc {
	return func(rw io.ReadWriter) (io.Writer, error) {
		target := &nio.BlockTarget{
			MetaWriter: rw,
			DataWriter: nio.WriterFunc(encode(cp, bobPub, rw)),
		}
		writer, err := nio.NewBlockTransfer(blockSize).Writer(target)
		if err != nil {
			return nil, fmt.Errorf("error sending block: %v", err)
		}
		return writer, nil
	}
}

func encode(
	cp dh.CommonPublicKey,
	bobPub *big.Int,
	rw io.ReadWriter,
) nio.WriterFunc {
	return func(p []byte) (int, error) {
		alice := elgamal_cipher.NewAlice(cp, bobPub)
		fmt.Print(fmtElgamalAlice(alice))

		msgInt := new(big.Int).SetBytes(p)
		log.Printf("ELGAMAL: data int: %v", msgInt)
		encoded := alice.Encode(msgInt, rand.CryptoSafeRandom())
		log.Printf("ELGAMAL: R=%v, E=%v", encoded.R, encoded.E)
		err := nio.WriteBigIntWithLen(rw, encoded.R)
		if err != nil {
			return 0, fmt.Errorf("writing R failed: %v", err)
		}
		err = nio.WriteBigIntWithLen(rw, encoded.E)
		if err != nil {
			return 0, fmt.Errorf("writing E failed: %v", err)
		}

		return len(p), nil
	}
}

// TODO: last byte is ignored
func ReceiveFunc(cp dh.CommonPublicKey) msg_core.ReceiveFunc {
	bob := elgamal_cipher.NewBob(cp)
	fmt.Print(fmtElgamalBob(bob))
	return func(rw io.ReadWriter) (io.Reader, error) {
		src := &nio.BlockSrc{
			MetaReader: rw,
			DataReader: nio.ReaderFunc(decrypt(bob, rw)),
		}
		return nio.NewBlockTransfer(blockSize).Reader(src), nil
	}
}

func decrypt(
	bob *elgamal_cipher.Bob,
	rw io.ReadWriter,
) nio.ReaderFunc {
	return func(buf []byte) (int, error) {
		R, err := nio.ReadBigIntWithLen(rw)
		if err != nil {
			return 0, fmt.Errorf("can't read R: %v", err)
		}

		E, err := nio.ReadBigIntWithLen(rw)
		if err != nil {
			return 0, fmt.Errorf("can't read E: %v", err)
		}
		encoded := &elgamal_cipher.Encoded{E: E, R: R}
		decoded := bob.Decode(encoded)
		if decoded.BitLen() > blockSize*8 {
			return 0, errors.New("received value is larger that buffer size. Seems like Alice uses incorrect key")
		}
		decoded.FillBytes(buf)
		log.Printf("ELGAMAL: decoded data=%v", buf)
		return blockSize, nil
	}
}

func fmtElgamalAlice(a *elgamal_cipher.Alice) string {
	return fmt.Sprintln("Elgamal node(Alice) initialized.\n",
		fmt.Sprintf("Common public key: P=%v, Q=%v\n", a.CommonPub.P(), a.CommonPub.G()),
		fmt.Sprintf("Bob public key: '%v'", a.BobPub),
	)
}

func fmtElgamalBob(bob *elgamal_cipher.Bob) string {
	err := ioutil.WriteFile(bobPubKeyFile, bob.Pub.Bytes(), 0644)
	if err != nil {
		return "error writing key"
	}
	return fmt.Sprintln("Elgamal node(Bob) initialized.\n",
		fmt.Sprintf("Common public key: P=%v, Q=%v\n", bob.CommonPub.P(), bob.CommonPub.G()),
		fmt.Sprintf("Node public key: '%v' (saved to %v)", bob.Pub, bobPubKeyFile),
	)
}
