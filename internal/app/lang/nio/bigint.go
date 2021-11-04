package nio

import (
	"fmt"
	"io"
	"math/big"
)

func WriteBigIntWithLen(writer io.Writer, msg *big.Int) error {
	msgBytes := msg.Bytes()

	msgLen := uint32(len(msgBytes))
	err := WriteUint32(writer, msgLen)
	if err != nil {
		return err
	}

	actuallyWrote, err := writer.Write(msgBytes)
	if err != nil {
		return err
	}
	// TODO: make it loop
	if uint32(actuallyWrote) != msgLen {
		return fmt.Errorf("tried to write %v byte but failed, written: %v", msgLen, actuallyWrote)
	}
	return nil
}

func ReadBigIntWithLen(reader io.Reader) (*big.Int, error) {
	msgLen, err := ReadUint32(reader)
	if err != nil {
		return nil, err
	}
	msgBuf := make([]byte, msgLen)
	actuallyRead, err := reader.Read(msgBuf)
	if err != nil {
		return nil, err
	}
	// TODO: make it loop
	if actuallyRead != int(msgLen) {
		return nil, fmt.Errorf("failed to read %v bytes, %v", msgLen, actuallyRead)
	}
	msg := new(big.Int).SetBytes(msgBuf)

	return msg, nil
}
