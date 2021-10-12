package common

import (
	. "fmt"
	"strings"
)

type CmdConf interface {
	CmdName() string
	InitCmd(args []string) (Cmd, CmdConfError)
}

type Cmd interface {
	Run() error
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
		subCmd, err := subConfig.InitCmd(args[1:])
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
	return NewCmdConfError(
		Sprintf("one of subcommands required: %v", strings.Join(subNames, ", ")),
		nil,
	)
}
