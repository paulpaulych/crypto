package nio

import (
	"fmt"
	"io"
	"log"
)

// BlockTransfer allows transferring data in parts(blocks).
// To achieve this, along with the original
//  data BlockTransfer writes/reads additional(i.e. meta) information.
// Current implementation writes/reads meta-byte before each block of original data,
// which means number of first corresponding block bytes that should be read (contain original data).
// So blockSize now cannot be greater than 255
type BlockTransfer struct {
	blockSize int
}

func NewBlockTransfer(blockSize int) *BlockTransfer {

	return &BlockTransfer{blockSize: blockSize}
}

type WriteProps struct {
	From       io.Reader
	MetaWriter io.Writer
	DataWriter io.Writer
}

func (b BlockTransfer) WriteBlocks(
	props WriteProps,
) error {
	buf := make([]byte, b.blockSize)
	for {
		actualRead, err := props.From.Read(buf)
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return fmt.Errorf("WriteBlocks: error reading block: %v", err)
		}

		{
			_, err = props.MetaWriter.Write([]byte{byte(actualRead)})
			if err != nil {
				return fmt.Errorf("WriteBlocks: writing metadata failed: %v", err)
			}
		}

		_, err = props.DataWriter.Write(buf)
		if err != nil {
			return fmt.Errorf("WriteBlocks: error writing block: %v", err)
		}
	}
}

type ReadProps struct {
	To         io.Writer
	MetaReader io.Reader
	DataReader io.Reader
}

func (b BlockTransfer) ReadBlocks(props ReadProps) error {
	metaBuf := []byte{1}
	metaByte := byte(0)
	blockBuf := make([]byte, b.blockSize)

	for {
		_, err := props.MetaReader.Read(metaBuf)
		if err != nil {
			return fmt.Errorf("error reading meta byte of block")
		}
		metaByte = metaBuf[0]

		_, err = props.DataReader.Read(blockBuf)
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return fmt.Errorf("can't read block: %v", err)
		}

		log.Printf("received block: len=%v, bytes=%v", len(blockBuf), blockBuf)
		_, err = props.To.Write(blockBuf[:int(metaByte)])
		if err != nil {
			return fmt.Errorf("error writing received message: %v", err)
		}
	}
}
