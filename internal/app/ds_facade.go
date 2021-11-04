package app

import (
	"fmt"
	"io/ioutil"
	"math/big"
	"os"

	"github.com/paulpaulych/crypto/internal/app/digital-sign"
)

const (
	signedFileName = "signed"
	secretFileName = "rsa.key"
 	publicFileName = "rsa_pub.key"

 	keyFilePerm = 0777
)

func RsaKeyGen(p, q *big.Int) error {
	keys, e := digital_sign.GenerateKeys(p, q)
	if e != nil {
		return nil
	}
	e = ioutil.WriteFile(secretFileName, keys.Secret, keyFilePerm)
	if e != nil {
		return fmt.Errorf("can't write secret file: %v", e)
	}
	e = ioutil.WriteFile(publicFileName, keys.Public, keyFilePerm)
	if e != nil {
		return fmt.Errorf("can't write public file: %v", e)
	}
	return nil
}

func RsaSign(msgFile string, secretKeyFile string) error {
	sec, e := ioutil.ReadFile(secretKeyFile)
	if e != nil {
		return fmt.Errorf("can't read secret file: %v", e)
	}

	msg, e := os.Open(msgFile)
	if e != nil {
		return nil
	}
	defer msg.Close()

	signed, e := digital_sign.Sign(msg, sec)
	if e != nil {
		return fmt.Errorf("signature error: %v", e)
	}

	e = ioutil.WriteFile(signedFileName, signed, keyFilePerm)
	if e != nil {
		return fmt.Errorf("can't write public file: %v", e)
	}
	return nil
}

func RsaValidate(signedFile string, pubFName string) error {
	signed, e := ioutil.ReadFile(signedFile)
	if e != nil {
		return fmt.Errorf("can't read signed file: %v", e)
	}
	pub, e := ioutil.ReadFile(pubFName)
	if e != nil {
		return fmt.Errorf("can't read public key: %v", e)
	}

	return digital_sign.Validate(signed, pub)	
}