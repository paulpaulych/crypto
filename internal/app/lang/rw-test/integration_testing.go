package wr_test

import (
	"io"
	"sync"
	"testing"
)

type ReadWriteTest struct {
	Name string
	A    Peer
	B    Peer
}

type Result interface{}

type Peer struct {
	Act   func(io.ReadWriter) (Result, error)
	Check func(got Result, gotE error) error
}

func RunReadWriteTest(t *testing.T, tt *ReadWriteTest) {
	conn1, conn2 := newPipe()

	t.Run(tt.Name, func(t *testing.T) {
		wg := sync.WaitGroup{}
		wg.Add(1)

		var (
			aRes interface{}
			aErr error
		)
		go func() {
			aRes, aErr = tt.A.Act(conn1)
			wg.Done()
		}()

		bRes, bErr := tt.B.Act(conn2)
		wg.Wait()

		aCheckErr := tt.A.Check(aRes, aErr)
		if aCheckErr != nil {
			t.Errorf("peer A test fail: %v", aCheckErr)
		}
		bCheckErr := tt.B.Check(bRes, bErr)
		if bCheckErr != nil {
			t.Errorf("peer B test fail: %v", bCheckErr)
		}
	})
}

type pipe struct {
	*io.PipeReader
	*io.PipeWriter
}

func newPipe() (io.ReadWriter, io.ReadWriter) {
	r1, w1 := io.Pipe()
	r2, w2 := io.Pipe()

	return &pipe{r1, w2}, &pipe{r2, w1}
}
