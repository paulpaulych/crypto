package tcp

import (
	"bytes"
	"encoding/binary"
	"errors"
	"net"
)

func WriteUint32(conn net.Conn, int uint32) error {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.BigEndian, int)
	if err != nil {
		return err
	}

	actuallyWrote, err := conn.Write(buf.Bytes())
	if err != nil {
		return err
	}
	if actuallyWrote != 4 {
		return errors.New("tried to write uint32 but failed")
	}
	return nil
}

func ReadUint32(conn net.Conn) (uint32, error) {
	lenBuf := make([]byte, 4)
	actuallyRead, err := conn.Read(lenBuf)
	if err != nil {
		return 0, err
	}
	if actuallyRead != 4 {
		return 0, errors.New("failed to read a byte")
	}
	return binary.BigEndian.Uint32(lenBuf), nil
}
