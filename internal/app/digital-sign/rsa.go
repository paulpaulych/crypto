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
	"github.com/paulpaulych/crypto/internal/core/rsa-ds"
)

// TODO: add format for saved files, checksums maybe

func Sign(msgReader io.Reader, secretKeyBytes []byte) ([]byte, error) {
	key, err := readSecret(secretKeyBytes)
	if err != nil {
		return nil, fmt.Errorf("error reading secret key from file: %v", err)
	}

	msg, err := ioutil.ReadAll(msgReader)
	if err != nil {
		return nil, err
	}

	signed, err := rsa_ds.Sign(
		key,
		&rsa_ds.Msg{Value: new(big.Int).SetBytes(msg)},
		simpleHash,
	)
	if err != nil {
		return nil, fmt.Errorf("signature error: %v", err)
	}

	signedBytes, err := writeSignature(signed)
	if err != nil {
		return nil, err
	}

	return signedBytes, nil
}

func Validate(
	msgBytes []byte,
	signBytes []byte,
	pubKeyBytes []byte,
) error {
	pubKey, err := readPublic(pubKeyBytes)
	if err != nil {
		return fmt.Errorf("error reading public key from file: %v", err)
	}
	msg:= new(big.Int).SetBytes(msgBytes)

	signed, err := readSignature(signBytes)
	if err != nil {
		return fmt.Errorf("error reading signed key from file: %v", err)
	}

	valid, err := rsa_ds.IsSignatureValid(
		pubKey,
		&rsa_ds.Msg{Value: msg},
		signed,
		simpleHash,
	)
	if err != nil {
		return fmt.Errorf("error validation signature: %v", err)
	}
	if !valid {
		fmt.Println("SIGNATURE IS INVALID")
	} else {
		fmt.Println("SIGNATURE IS VALID")
	}
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

func readSignature(from []byte) (*rsa_ds.Signature, error) {
	buf := bytes.NewBuffer(from)

	signature, err := nio.ReadBigIntWithLen(buf)
	if err != nil {
		return nil, err
	}
	return &rsa_ds.Signature{Value: signature}, nil
}

func writeSignature(sign *rsa_ds.Signature) ([]byte, error) {
	buf := &bytes.Buffer{}
	err := nio.WriteBigIntWithLen(buf, sign.Value)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
