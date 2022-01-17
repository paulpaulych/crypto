package app

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net"
	"os"

	"github.com/paulpaulych/crypto/internal/app/lang/nio"
	"github.com/paulpaulych/crypto/internal/app/lang/tcp"
	"github.com/paulpaulych/crypto/internal/app/messaging/msg-core"
	"github.com/paulpaulych/crypto/internal/app/messaging/protocols"
)

func ShamirRecv(bindAddr string) error {
	return tcp.StartServer(
		bindAddr,
		loggingErrorsChan(),
		func(conn net.Conn) error {
			return decryptAndSaveToFile(conn, protocols.ShamirReceiver())
		},
	)
}

func ShamirSend(target string, fname string, p *big.Int) error {
	protocol, e := protocols.ShamirSender(p)
	if e != nil {
		return fmt.Errorf("protocol error: %v", e)
	}
	return sendEncryptedFile(target, fname, protocol)
}

func ElgamalRecv(bindAddr string, p, g *big.Int) error {
	protocol, e := protocols.ElgamalReceiver(p, g)
	if e != nil {
		return fmt.Errorf("protocol error: %v", e)
	}
	// TODO: log result
	return tcp.StartServer(
		bindAddr,
		loggingErrorsChan(),
		func(conn net.Conn) error {
			return decryptAndSaveToFile(conn, protocol)
		},
	)
}

func ElgamalSend(
	target string,
	fname string,
	p, g *big.Int,
	bobPubKeyFile string,
) error {
	protocol, e := protocols.ElgamalSender(p, g, bobPubKeyFile)
	if e != nil {
		return fmt.Errorf("protocol error: %v", e)
	}
	// TODO: log result 
	return sendEncryptedFile(target, fname, protocol)
}

func RsaRecv(bindAddr string, p, q *big.Int) error {
	receiver, err:= protocols.RsaReceiver(p, q)
	if err != nil {
		return err
	}
	return tcp.StartServer(
		bindAddr,
		loggingErrorsChan(),
		func(conn net.Conn) error {
			return decryptAndSaveToFile(conn, receiver)
		},
	)
}

func RsaSend(target string, fname string, bobPubKeyFile string) error {
	bobKey, e:= ioutil.ReadFile(bobPubKeyFile)
	if e != nil {
		return fmt.Errorf("can't read bob key: %v", e)
	}
	sender, e := protocols.RsaSender(bobKey)
	if e != nil {
		return e
	}
	return sendEncryptedFile(target, fname, sender)
}

func loggingErrorsChan() chan<- error {
	ee := make(chan error)

	logger := func() {
		for e := range ee {
			log.Printf("ERROR: %v", e)
		}
	}

	go logger()

	return ee
}

func sendEncryptedFile(
	target string,
	fname string,
	protocol msg_core.Sender,
) error {
	f, e := os.Open(fname)
	if e != nil {
		return e
	}
	fileData := nio.NewFileData(f)
	return sendEncryptedMessage(target, fileData, protocol)
}

func decryptAndSaveToFile(conn net.Conn, protocol msg_core.Receiver) error {
	defer conn.Close()
	reader, e := msg_core.EncryptedMsgReader(conn, protocol)
	if e != nil {
		return fmt.Errorf("error reading encrypted msg: %v", e)
	}
	fileData, e := nio.ReadFileData(reader)
	if e != nil {
		return fmt.Errorf("error reading file: %v", e)
	}
	e = fileData.WriteToFile()
	if e != nil {
		return fmt.Errorf("error writing file: %v", e)
	}
	return nil
}


func sendEncryptedMessage(addr string, data *nio.FileData, protocol msg_core.Sender) error {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return fmt.Errorf("can't connect to %s: %v", addr, err)
	}
	defer conn.Close()
	log.Printf("connected to %s", addr)
	wr, err := msg_core.EncryptedMsgWriter(conn, protocol)
	if err != nil {
		return fmt.Errorf("can't")
	}
	_, err = data.WriteTo(wr)
	if err != nil {
		return err
	}
	return nil
}
