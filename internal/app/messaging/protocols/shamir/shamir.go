package shamir

import (
	"fmt"
	"io"
	"math/big"
	"github.com/paulpaulych/crypto/internal/app/lang/nio"
	"github.com/paulpaulych/crypto/internal/app/messaging/msg-core"
	"github.com/paulpaulych/crypto/internal/core/shamir-cipher"
)

func calcBlockSize(p *big.Int) uint {
	return uint(len(p.Bytes()) - 1)
}

const minP = 256

func SendFunc(p *big.Int) (msg_core.SendFunc, error) {
	if p.Cmp(big.NewInt(minP)) <= 0 {
		return nil, fmt.Errorf("p must be greater than %v", minP)
	}
	sendFunc := func(rw io.ReadWriter) (io.Writer, error) {
		err := nio.WriteBigIntWithLen(rw, p)
		if err != nil {
			return nil, fmt.Errorf("writing P failed: %v", err)
		}
		blockSize := calcBlockSize(p)

		target := &nio.BlockTarget {
			MetaWriter: rw,
			DataWriter: encrypt(p, rw),
		}
		writer, err := nio.NewBlockTransfer(blockSize).Writer(target)
		if err != nil {
			return nil, fmt.Errorf("error sending block: %v", err)
		}
		return writer, nil
	}
	return sendFunc, nil
}

func encrypt(p *big.Int, rw io.ReadWriter) nio.WriterFunc {
	return func(buf []byte) (int, error) {
		alice, err := shamir_cipher.InitAlice(p)
		if err != nil {
			return 0, fmt.Errorf("failed to init alice: %d", err)
		}
	
		step1out, err := alice.Step1(new(big.Int).SetBytes(buf))
		if err != nil {
			return 0, fmt.Errorf("writing step1out failed: %v", err)
		}
	
		err = nio.WriteBigIntWithLen(rw, step1out)
		if err != nil {
			return 0, fmt.Errorf("writing step1out failed: %v", err)
		}
	
		step2out, err := nio.ReadBigIntWithLen(rw)
		if err != nil {
			return 0, fmt.Errorf("reading step2out failed: %v", err)
		}
	
		step3out := alice.Step3(step2out)
		err = nio.WriteBigIntWithLen(rw, step3out)
		if err != nil {
			return 0, fmt.Errorf("writing step3out failed: %v", err)
		}
	
		return 0, nil
	}
}

func ReceiveFunc(rw io.ReadWriter) (io.Reader, error) {
	p, err := nio.ReadBigIntWithLen(rw)
	if err != nil {
		return nil, fmt.Errorf("can't read p: %s", err)
	}

	blockSize := calcBlockSize(p)
	opts := &nio.BlockSrc{
		MetaReader: rw,
		DataReader: decrypt(p, rw),
	}
	reader := nio.NewBlockTransfer(blockSize).Reader(opts)
	return reader, nil
}

func decrypt(p *big.Int, rw io.ReadWriter) nio.ReaderFunc {
	return func(buf []byte) (int, error) {
		bob, err := shamir_cipher.InitBob(p)
		if err != nil {
			return 0, fmt.Errorf("failed to init bob: %d", err)
		}
	
		step1out, err := nio.ReadBigIntWithLen(rw)
		if err == io.EOF {
			return 0, io.EOF
		}
		if err != nil {
			return 0, fmt.Errorf("can't read step1out: %v", err)
		}
	
		step2out := bob.Step2(step1out)
		err = nio.WriteBigIntWithLen(rw, step2out)
		if err != nil {
			return 0, fmt.Errorf("can't write step2out: %v", err)
		}
	
		step3out, err := nio.ReadBigIntWithLen(rw)
		if err != nil {
			return 0, fmt.Errorf("can't write step2out: %v", err)
		}
		bob.Decode(step3out).FillBytes(buf)
	
		return len(buf), nil
	}	
}
