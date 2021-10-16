package tcp

import (
	"fmt"
	"log"
	. "math/big"
	"net"
)

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
		return fmt.Errorf("tried to write %v byte but failed, written: %v", msgLen, actuallyWrote)
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
		return nil, fmt.Errorf("failed to read %v bytes, %v", msgLen, actuallyRead)
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
