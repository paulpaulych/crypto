package nio

type ClosableWriter interface {
	Write(p []byte) (int, error)
	Close() error
}
