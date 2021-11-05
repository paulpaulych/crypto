package app

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"os"

	"github.com/paulpaulych/crypto/internal/app/digital-sign"
)

const (
	signatureFileName = "signature"
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
	log.Println("SECRET KEY SAVED TO", secretFileName)
	e = ioutil.WriteFile(publicFileName, keys.Public, keyFilePerm)
	if e != nil {
		return fmt.Errorf("can't write public file: %v", e)
	}
	log.Println("PUBLIC KEY SAVED TO", publicFileName)
	return nil
}

func RsaSign(msgFile string, secretKeyFile string) error {
	sec, e := ioutil.ReadFile(secretKeyFile)
	if e != nil {
		return fmt.Errorf("can't read secret file: %v", e)
	}

	msg, e := os.Open(msgFile)
	if e != nil {
		return fmt.Errorf("can't open msg file: %v", e)
	}
	defer msg.Close()

	sign, e := digital_sign.Sign(msg, sec)
	if e != nil {
		return fmt.Errorf("signature error: %v", e)
	}
	log.Println("SIGNATURE SAVED TO", signatureFileName)

	e = ioutil.WriteFile(signatureFileName, sign, keyFilePerm)
	if e != nil {
		return fmt.Errorf("can't write public file: %v", e)
	}
	return nil
}

func RsaValidate(msgFile string, signFile string, pubFName string) error {
	msg, e := ioutil.ReadFile(msgFile)
	if e != nil {
		return fmt.Errorf("can't open msg file: %v", e)
	}
	signature, e := ioutil.ReadFile(signFile)
	if e != nil {
		return fmt.Errorf("can't read signed file: %v", e)
	}
	pub, e := ioutil.ReadFile(pubFName)
	if e != nil {
		return fmt.Errorf("can't read public key: %v", e)
	}

	return digital_sign.Validate(msg, signature, pub)	
}