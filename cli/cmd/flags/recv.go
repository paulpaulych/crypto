package flags

import (
	"flag"
	"fmt"
	cmd "github.com/paulpaulych/crypto/shamir/cli/cmd"
)

func ParseRecvCmdFlags(flags *flag.FlagSet, args []string) (*cmd.RecvArgs, error) {
	bindHostPtr := flags.String("host", "localhost", "host to bind")
	bindPortPtr := flags.Int("port", 0, "port to bind")

	err := flags.Parse(args)
	if err != nil {
		return nil, err
	}

	addr := fmt.Sprintf("%s:%v", *bindHostPtr, *bindPortPtr)
	return &cmd.RecvArgs{BindAddr: addr}, nil
}
