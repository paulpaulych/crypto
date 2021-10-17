package tcp

import (
	"fmt"
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

	return msg, nil
}
