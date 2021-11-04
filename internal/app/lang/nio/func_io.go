package nio

type ReaderFunc func(p []byte) (int, error)

func (f ReaderFunc) Read(p []byte) (int, error) {
	return f(p)
}

type WriterFunc func(p []byte) (int,error)

func (f WriterFunc) Write(p []byte) (int, error) {
	return f(p)
}
