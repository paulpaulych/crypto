package nio

type BlockWriter interface {
	Write(p []byte, hasMore bool) error
}

type BytePage struct {
	Bytes   []byte
	HasMore bool
}

type ByteReader interface {
	Read(n uint) (p *BytePage, err error)
}

func ReadAll(reader ByteReader) ([]byte, error) {
	res := make([]byte, 0)
	blockSize := uint(1)
	for {
		page, err := reader.Read(blockSize)
		if err != nil {
			return nil, err
		}
		res = append(res, page.Bytes...)
		if !page.HasMore {
			return res, nil
		}
	}
}
