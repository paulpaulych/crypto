package digital_sign

import (
	"fmt"
	"github.com/paulpaulych/crypto/internal/app/nio"
	"github.com/paulpaulych/crypto/internal/core/rand"
	"github.com/paulpaulych/crypto/internal/core/rsa-ds"
	"io"
	"io/ioutil"
	. "math/big"
	"os"
)

const signedFileName = "signed"

func Sign(msgReader io.Reader, secretKeyFile string) error {
	keyF, err := os.Open(secretKeyFile)
	if err != nil {
		return err
	}
	defer func() { _ = keyF.Close() }()

	key, err := readSecret(keyF)
	if err != nil {
		return fmt.Errorf("error reading secret key from file: %v", err)
	}

	msg, err := ioutil.ReadAll(msgReader)
	if err != nil {
		return err
	}

	signed, err := rsa_ds.Sign(key, new(Int).SetBytes(msg), simpleHash)
	if err != nil {
		return fmt.Errorf("signature error: %v", err)
	}

	signedFile := nio.NewFileWriter(signedFileName, func() {
		fmt.Println("SIGNED MESSAGE SAVED TO", signedFileName)
	})
	defer func() { _ = signedFile.Close() }()
	err = writeSigned(signedFile, signed)
	if err != nil {
		return fmt.Errorf("error writing signed message to file: %v", err)
	}

	return nil
}

func Validate(signedFile string, pubKeyFile string) error {
	pubKeyF, err := os.Open(pubKeyFile)
	if err != nil {
		return err
	}
	defer func() { _ = pubKeyF.Close() }()

	pubKey, err := readPublic(pubKeyF)
	if err != nil {
		return fmt.Errorf("error reading public key from file: %v", err)
	}

	signedF, err := os.Open(signedFile)
	if err != nil {
		return err
	}
	defer func() { _ = signedF.Close() }()

	signed, err := readSigned(signedF)
	if err != nil {
		return fmt.Errorf("error reading public key from file: %v", err)
	}

	valid, err := rsa_ds.IsSignatureValid(pubKey, signed, simpleHash)
	if err != nil {
		return fmt.Errorf("error validation signature: %v", err)
	}
	if !valid {
		fmt.Println("SIGNATURE IS INVALID")
	} else {
		fmt.Println("SIGNATURE IS VALID")
	}
	fmt.Printf("Message: %v\n", string(signed.Msg.Bytes()))
	return nil
}

const secretFileName = "rsa.key"
const publicFileName = "rsa_pub.key"

func GenerateKeys(p, q *Int) error {
	pub, sec, err := rsa_ds.GenKeys(p, q, rand.CryptoSafeRandom())
	if err != nil {
		return fmt.Errorf("error generationg keys: %v", err)
	}
	secretFile := nio.NewFileWriter(secretFileName, func() {
		fmt.Println("SECRET KEY SAVED TO", secretFileName)
	})
	defer func() { _ = secretFile.Close() }()
	err = writeSecret(sec, secretFile)
	if err != nil {
		return fmt.Errorf("error writing secret key to file: %v", err)
	}

	publicFile := nio.NewFileWriter(publicFileName, func() {
		fmt.Println("PUBLIC KEY SAVED TO", publicFileName)
	})
	defer func() { _ = publicFile.Close() }()
	err = writePublic(pub, publicFile)
	if err != nil {
		return fmt.Errorf("error writing public key to file: %v", err)
	}
	return nil
}

func writePublic(key *rsa_ds.PubKey, writer io.Writer) error {
	err := nio.WriteBigIntWithLen(writer, key.N)
	if err != nil {
		return err
	}
	err = nio.WriteBigIntWithLen(writer, key.Exp)
	if err != nil {
		return err
	}
	return nil
}

func readPublic(reader io.Reader) (*rsa_ds.PubKey, error) {
	N, err := nio.ReadBigIntWithLen(reader)
	if err != nil {
		return nil, err
	}
	Exp, err := nio.ReadBigIntWithLen(reader)
	if err != nil {
		return nil, err
	}
	return &rsa_ds.PubKey{Exp: Exp, N: N}, nil
}

func writeSecret(key *rsa_ds.SecretKey, writer io.Writer) error {
	err := nio.WriteBigIntWithLen(writer, key.N)
	if err != nil {
		return err
	}
	err = nio.WriteBigIntWithLen(writer, key.Exp)
	if err != nil {
		return err
	}
	return nil
}

func readSecret(reader io.Reader) (*rsa_ds.SecretKey, error) {
	N, err := nio.ReadBigIntWithLen(reader)
	if err != nil {
		return nil, err
	}
	Exp, err := nio.ReadBigIntWithLen(reader)
	if err != nil {
		return nil, err
	}
	return &rsa_ds.SecretKey{Exp: Exp, N: N}, nil
}

func simpleHash(orig *Int) (*Int, error) {
	return new(Int).SetBytes(orig.Bytes()[:1]), nil
}

func readSigned(reader io.Reader) (*rsa_ds.Signed, error) {
	msg, err := nio.ReadBigIntWithLen(reader)
	if err != nil {
		return nil, err
	}
	signature, err := nio.ReadBigIntWithLen(reader)
	if err != nil {
		return nil, err
	}
	return &rsa_ds.Signed{Msg: msg, Signature: signature}, nil
}

func writeSigned(writer io.Writer, signed *rsa_ds.Signed) error {
	err := nio.WriteBigIntWithLen(writer, signed.Msg)
	if err != nil {
		return err
	}
	err = nio.WriteBigIntWithLen(writer, signed.Signature)
	if err != nil {
		return err
	}
	return nil
}
