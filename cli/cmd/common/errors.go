package common

type CmdConfError interface {
	Error() string
	Trace() []string
	Usage() *Usage
}

type confError struct {
	trace []string
	error string
	usage *Usage
}

func (e *confError) Error() string {
	return e.error
}
func (e *confError) Trace() []string {
	return e.trace
}
func (e *confError) Usage() *Usage {
	return e.usage
}

type Usage = string

func NewCmdConfError(msg string, usage *Usage) CmdConfError {
	return &confError{
		trace: []string{},
		error: msg,
		usage: usage,
	}
}

func AppendTrace(e CmdConfError, subCmdName string) CmdConfError {
	return &confError{
		trace: append([]string{subCmdName}, e.Trace()...),
		error: e.Error(),
		usage: e.Usage(),
	}
}
