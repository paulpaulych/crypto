package cli

import (
	"fmt"
	"io"
	"strings"
)

type Cmd interface {
	Run() error
}

type CmdConf interface {
	CmdName() string
	NewCmd(args []string) (Cmd, CmdConfError)
}

type WriteHelp = func(writer io.Writer)

type CmdConfError interface {
	Error() error
	Trace() []string
	HelpWriter() WriteHelp
}

func InitSubCmd(subConfigs []CmdConf, args []string) (Cmd, CmdConfError) {
	if len(args) < 1 {
		return nil, noSubCmdError(subConfigs)
	}

	subCmdName := args[0]

	for _, subConfig := range subConfigs {
		if subCmdName != subConfig.CmdName() {
			continue
		}
		subCmd, err := subConfig.NewCmd(args[1:])
		if err != nil {
			return nil, AppendTrace(err, subConfig.CmdName())
		}
		return subCmd, nil
	}

	return nil, noSubCmdError(subConfigs)
}

func noSubCmdError(configs []CmdConf) CmdConfError {
	subNames := make([]string, len(configs))
	for i, subConfig := range configs {
		subNames[i] = subConfig.CmdName()
	}
	return NewCmdConfErr(
		fmt.Errorf("one of subcommands required: %v", strings.Join(subNames, ", ")),
		nil,
	)
}
