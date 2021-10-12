package nio

type BlockWriter interface {
	Write(p []byte, hasMore bool) error
}

type BytePage struct {
	Bytes   []byte
	HasMore bool
}

type ByteReader interface {
	Read(n uint) (*BytePage, error)
	TotalBytes() (uint, error)
}
