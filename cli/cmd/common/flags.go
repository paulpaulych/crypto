package common

import (
	"flag"
	"strings"
)

func Parse(fs *flag.FlagSet, args []string) CmdConfError {
	usageBuilder := &strings.Builder{}
	fs.SetOutput(usageBuilder)
	err := fs.Parse(args)
	if err == nil {
		return nil
	}
	usage := usageBuilder.String()
	return NewCmdConfError(err.Error(), &usage)
}
