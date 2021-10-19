package digital_sign

import (
	"fmt"
	"github.com/paulpaulych/crypto/internal/app/nio"
	"github.com/paulpaulych/crypto/internal/core/rand"
	"github.com/paulpaulych/crypto/internal/core/rsa-ds"
	"io"
	. "math/big"
)

func Sign(msgFile string, privateKeyFile string) {

}

func Validate(msgFile string, pubKeyFile string) {

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
