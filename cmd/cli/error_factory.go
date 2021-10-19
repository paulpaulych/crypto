package cli

func NewCmdConfErr(err error, usage WriteHelp) CmdConfError {
	return &confError{
		trace: []string{},
		error: err,
		usage: usage,
	}
}

func HelpRequested(usage WriteHelp) CmdConfError {
	return &confError{
		trace: []string{},
		error: nil,
		usage: usage,
	}
}

func AppendTrace(e CmdConfError, subCmdName string) CmdConfError {
	return &confError{
		trace: append([]string{subCmdName}, e.Trace()...),
		error: e.Error(),
		usage: e.HelpWriter(),
	}
}

type confError struct {
	trace []string
	error error
	usage WriteHelp
}

func (e *confError) Error() error {
	return e.error
}
func (e *confError) Trace() []string {
	return e.trace
}
func (e *confError) HelpWriter() WriteHelp {
	return e.usage
}
