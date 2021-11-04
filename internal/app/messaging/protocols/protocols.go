package protocols

import (
	"fmt"
	"io/ioutil"
	"math/big"

	"github.com/paulpaulych/crypto/internal/app/messaging/msg-core"
	"github.com/paulpaulych/crypto/internal/app/messaging/protocols/elgamal"
	"github.com/paulpaulych/crypto/internal/app/messaging/protocols/rsa"
	"github.com/paulpaulych/crypto/internal/app/messaging/protocols/shamir"
	dh "github.com/paulpaulych/crypto/internal/core/diffie-hellman"
)

const (
	Shamir msg_core.ProtocolCode = iota
	Elgamal
	Rsa
)

func ShamirReceiver() msg_core.Receiver {
	return msg_core.NewReceiver(Shamir, shamir.ReceiveFunc)
}

func ShamirSender(p *big.Int) (msg_core.Sender, error) {
	write, err := shamir.SendFunc(p)
	if err != nil {
		return nil, err
	}
	return msg_core.NewSender(Shamir, write), nil
}

func RsaSender(bobPub []byte) (msg_core.Sender, error) {
	fn, err := rsa.SendFunc(bobPub)
	if err != nil {
		return nil, err
	}
	return msg_core.NewSender(Rsa, fn), nil
}

func RsaReceiver(p, q *big.Int) (msg_core.Receiver, error) {
	receiveFunc, err := rsa.ReceiveFunc(p, q, saveRsaBobPubKeyToFile)
	if err != nil {
		return nil, err
	}
	return msg_core.NewReceiver(Rsa, receiveFunc), nil
}

func ElgamalSender(p, g *big.Int, bobPubFileName string) (msg_core.Sender, error) {
	commonPub, e := dh.NewCommonPublicKey(p, g)
	if e != nil {
		return nil, fmt.Errorf("Diffie-Hellman public key error: %v", e)
	}
	bytes, err := ioutil.ReadFile(bobPubFileName)
	if err != nil {
		return nil, err
	}
	bobPub := new(big.Int).SetBytes(bytes)
	writeFn := elgamal.SendFunc(commonPub, bobPub)
	return msg_core.NewSender(Elgamal, writeFn), nil
}

func ElgamalReceiver(p, g *big.Int) (msg_core.Receiver, error) {
	commonPub, e := dh.NewCommonPublicKey(p, g)
	if e != nil {
		return nil, fmt.Errorf("Diffie-Hellman public key error: %v", e)
	}
	return msg_core.NewReceiver(Elgamal, elgamal.ReceiveFunc(commonPub)), nil
}

const rsaBobKeyFileName = "bob_rsa.key"
func saveRsaBobPubKeyToFile(bobPub []byte) error {
	err := ioutil.WriteFile(rsaBobKeyFileName, bobPub, 0777)
	if err != nil {
		return err
	}
	fmt.Println("PUBLIC KEY SAVED TO", rsaBobKeyFileName)
	return nil
}
