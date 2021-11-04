package digital_sign

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/big"

	"github.com/paulpaulych/crypto/internal/app/lang/nio"
	"github.com/paulpaulych/crypto/internal/core/rand"
	rsa_ds "github.com/paulpaulych/crypto/internal/core/rsa-ds"
)

func Sign(msgReader io.Reader, secretKeyBytes []byte) ([]byte, error) {
	key, err := readSecret(secretKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("error reading secret key from file: %v", err)
	}

	msg, err := ioutil.ReadAll(msgReader)
	if err != nil {
		return nil, err
	}

	signed, err := rsa_ds.Sign(key, new(big.Int).SetBytes(msg), simpleHash)
	if err != nil {
		return nil, fmt.Errorf("signature error: %v", err)
	}

	signedBytes, err := writeSigned(signed)
	if err != nil {
		return nil, err
	}

	return signedBytes, nil
}

func Validate(signedReader []byte, pubKeyBytes []byte) error {
	pubKey, err := readPublic(pubKeyBytes)
	if err != nil {
		return fmt.Errorf("error reading public key from file: %v", err)
	}

	signed, err := readSigned(signedReader)
	if err != nil {
		return fmt.Errorf("error reading signed key from file: %v", err)
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

type Keys struct {
	Public []byte
	Secret []byte
}

func GenerateKeys(p, q *big.Int) (*Keys, error) {
	pub, sec, err := rsa_ds.GenKeys(p, q, rand.CryptoSafeRandom())
	if err != nil {
		return nil, fmt.Errorf("error generationg keys: %v", err)
	}

	secret, err := writeSecret(sec)
	if err != nil {
		log.Fatalf("error writing secret key to file: %v", err)
	}

	public, e := writePublic(pub)
	if e != nil {
		log.Fatalf("error writing public key to file: %v", err)
	}

	return &Keys{Public: public, Secret: secret}, nil
}

func writePublic(key *rsa_ds.PubKey) ([]byte, error) {
	buf := &bytes.Buffer{}
	err := nio.WriteBigIntWithLen(buf, key.N)
	if err != nil {
		return nil, err
	}
	err = nio.WriteBigIntWithLen(buf, key.Exp)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func readPublic(key []byte) (*rsa_ds.PubKey, error) {
	buf := bytes.NewBuffer(key)
	N, err := nio.ReadBigIntWithLen(buf)
	if err != nil {
		return nil, err
	}
	Exp, err := nio.ReadBigIntWithLen(buf)
	if err != nil {
		return nil, err
	}
	return &rsa_ds.PubKey{Exp: Exp, N: N}, nil
}

func writeSecret(key *rsa_ds.SecretKey) ([]byte, error) {
	buf := &bytes.Buffer{}
	err := nio.WriteBigIntWithLen(buf, key.N)
	if err != nil {
		return nil, err
	}
	err = nio.WriteBigIntWithLen(buf, key.Exp)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func readSecret(key []byte) (*rsa_ds.SecretKey, error) {
	buf := bytes.NewBuffer(key)
	N, err := nio.ReadBigIntWithLen(buf)
	if err != nil {
		return nil, err
	}
	Exp, err := nio.ReadBigIntWithLen(buf)
	if err != nil {
		return nil, err
	}
	return &rsa_ds.SecretKey{Exp: Exp, N: N}, nil
}

// TODO: use more copmplicated algorithm
func simpleHash(orig *big.Int) (*big.Int, error) {
	return new(big.Int).SetBytes(orig.Bytes()[:1]), nil
}

func readSigned(from []byte) (*rsa_ds.Signed, error) {
	buf := bytes.NewBuffer(from)
	msg, err := nio.ReadBigIntWithLen(buf)
	if err != nil {
		return nil, err
	}
	signature, err := nio.ReadBigIntWithLen(buf)
	if err != nil {
		return nil, err
	}
	return &rsa_ds.Signed{Msg: msg, Signature: signature}, nil
}

func writeSigned(signed *rsa_ds.Signed) ([]byte, error) {
	buf := &bytes.Buffer{}
	err := nio.WriteBigIntWithLen(buf, signed.Msg)
	if err != nil {
		return nil, err
	}
	err = nio.WriteBigIntWithLen(buf, signed.Signature)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
