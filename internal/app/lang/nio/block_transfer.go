package nio

import (
	"fmt"
	"io"
	"log"
)

const maxBlockSize = 255

// BlockTransfer allows transferring data in parts(blocks).
// To achieve this, along with the original
//  data  writes/reads meta-byte before each block of original data,
// which means number of first corresponding block bytes that should be read (contain original data).
// So blockSiBlockTransfer writes/reads additional(i.e. meta) information.
// Current implementationze now cannot be greater than 255
type BlockTransfer struct {
	blockSize uint
}

func NewBlockTransfer(blockSize uint) *BlockTransfer {
	if blockSize > maxBlockSize {
		log.Panicln("BlockTransfer: blockSize cannot be grater than", maxBlockSize)
	}
	log.Printf("BLOCK: Initialized. Block size = %v", blockSize)
	return &BlockTransfer{blockSize: blockSize}
}

type BlockTarget struct {
	MetaWriter io.Writer
	DataWriter io.Writer
}

func (b BlockTransfer) Writer(t *BlockTarget) (io.Writer, error) {
	block := make([]byte, b.blockSize)

	pr, pw := io.Pipe()

	go func() {
		for {
			actualRead, err := pr.Read(block)
			if err == io.EOF {
				return
			}
			if err != nil {
				fmt.Printf("error reading block: %v", err)
				return
			}
			log.Println("BLOCK: writing meta-byte:", actualRead)

			_, err = t.MetaWriter.Write([]byte{byte(actualRead)})
			if err != nil {
				fmt.Printf("writing metadata failed: %v", err)
				return
			}

			log.Println("BLOCK: writing data:", block)
			_, err = t.DataWriter.Write(block)
			if err != nil {
				fmt.Printf("error writing block: %v", err)
				return
			}
		}
	}()

	return pw, nil
}

type BlockSrc struct {
	MetaReader io.Reader
	DataReader io.Reader
}

func (b *BlockTransfer) Reader(src *BlockSrc) io.Reader {
	metaBuf := []byte{1}
	metaByte := byte(0)
	block := make([]byte, b.blockSize)

	pr, pw := io.Pipe()

	go func() {
		for {
			_, err := src.MetaReader.Read(metaBuf)
			if err == io.EOF {
				return
			}
			if err != nil {
				fmt.Printf("error reading meta byte of block: %v", err)
				return
			}
			metaByte = metaBuf[0]
			log.Println("BLOCK: received meta-byte:", metaByte)

			_, err = src.DataReader.Read(block)
			if err != nil {
				fmt.Printf("can't read block: %v", err)
				return
			}

			log.Println("BLOCK: received data:", block)
			_, err = pw.Write(block[:int(metaByte)])
			if err != nil {
				fmt.Printf("error writing received message: %v", err)
				return
			}
		}
	}()
	return pr
}
