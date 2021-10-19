package cli

import (
	"github.com/jessevdk/go-flags"
)

func ParseFlagsOfCmd(cmdName string, opts interface{}, args []string) ([]string, CmdConfError) {
	parser := flags.NewNamedParser(cmdName, flags.HelpFlag)
	parser.ArgsRequired = true
	parser.Args()
	_, err := parser.AddGroup("Command options", "", opts)
	if err != nil {
		return nil, NewCmdConfErr(err, nil)
	}
	remainingArgs, err := parser.ParseArgs(args)
	if err == nil {
		return remainingArgs, nil
	}
	switch flagsErr := err.(type) {
	case flags.ErrorType:
		if flagsErr == flags.ErrHelp {
			return remainingArgs, HelpRequested(nil)
		}
		return remainingArgs, NewCmdConfErr(err, nil)
	default:
		return remainingArgs, NewCmdConfErr(err, nil)
	}
}

type ParsedSubCmd struct {
	subCmdName    string
	remainingArgs []string
}
