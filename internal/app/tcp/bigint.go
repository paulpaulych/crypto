package tcp

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"log"
	. "math/big"
	"net"
)

const maxMsgLen = 255

func WriteBigIntWithLen(conn net.Conn, msg *Int) error {
	msgBytes := msg.Bytes()

	msgLen := len(msgBytes)
	if msgLen > maxMsgLen {
		err := fmt.Sprintf("cannot send cli larger than %v", maxMsgLen)
		return errors.New(err)
	}
	lenBuf := new(bytes.Buffer)
	err := binary.Write(lenBuf, binary.BigEndian, byte(msgLen))
	if err != nil {
		return err
	}

	logWriting(msg, lenBuf.Bytes(), msgBytes)

	actuallyWrote, err := conn.Write([]byte{byte(msgLen)})
	if err != nil {
		return err
	}
	if actuallyWrote != 1 {
		return errors.New("tried to write byte but failed")
	}

	actuallyWrote, err = conn.Write(msgBytes)
	if err != nil {
		return err
	}
	if actuallyWrote != msgLen {
		errMsg := fmt.Sprintf("tried to write %v byte but failed, written: %v", msgLen, actuallyWrote)
		return errors.New(errMsg)
	}
	return nil
}

func ReadBigIntWithLen(conn net.Conn) (*Int, error) {
	lenBuf := make([]byte, 1)
	actuallyRead, err := conn.Read(lenBuf)
	if err != nil {
		return nil, err
	}
	if actuallyRead != 1 {
		return nil, errors.New("failed to read a byte")
	}

	msgBuf := make([]byte, lenBuf[0])
	actuallyRead, err = conn.Read(msgBuf)
	if err != nil {
		return nil, err
	}
	if actuallyRead != int(lenBuf[0]) {
		errMsg := fmt.Sprintf("failed to read %v bytes, %v", lenBuf[0], actuallyRead)
		return nil, errors.New(errMsg)
	}
	msg := new(Int).SetBytes(msgBuf)

	logRead(msg, lenBuf, msgBuf)

	return msg, nil
}

func logRead(msg *Int, size, bytes []byte) {
	log.Printf("TCP: received Int %v: len=%v, bytes=%v", msg, size, bytes)
}

func logWriting(msg *Int, size, bytes []byte) {
	log.Printf("TCP: sending Int %v: len=%v, bytes=%v", msg, size, bytes)
}
