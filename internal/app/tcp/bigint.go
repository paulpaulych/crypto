package tcp

import (
	"errors"
	"fmt"
	"log"
	. "math/big"
	"net"
)

const maxMsgLen = 255

func WriteBigIntWithLen(conn net.Conn, msg *Int) error {
	msgBytes := msg.Bytes()

	msgLen := uint32(len(msgBytes))
	err := WriteUint32(conn, msgLen)
	if err != nil {
		return err
	}
	logWriting(msg, msgLen, msgBytes)

	actuallyWrote, err := conn.Write(msgBytes)
	if err != nil {
		return err
	}
	if uint32(actuallyWrote) != msgLen {
		errMsg := fmt.Sprintf("tried to write %v byte but failed, written: %v", msgLen, actuallyWrote)
		return errors.New(errMsg)
	}
	return nil
}

func ReadBigIntWithLen(conn net.Conn) (*Int, error) {
	msgLen, err := ReadUint32(conn)
	if err != nil {
		return nil, err
	}
	msgBuf := make([]byte, msgLen)
	actuallyRead, err := conn.Read(msgBuf)
	if err != nil {
		return nil, err
	}
	if actuallyRead != int(msgLen) {
		errMsg := fmt.Sprintf("failed to read %v bytes, %v", msgLen, actuallyRead)
		return nil, errors.New(errMsg)
	}
	msg := new(Int).SetBytes(msgBuf)

	logRead(msg, msgLen, msgBuf)

	return msg, nil
}

func logRead(msg *Int, size uint32, bytes []byte) {
	log.Printf("TCP: received Int %v: len=%v, bytes=%v", msg, size, bytes)
}

func logWriting(msg *Int, size uint32, bytes []byte) {
	log.Printf("TCP: sending Int %v: len=%v, bytes=%v", msg, size, bytes)
}
