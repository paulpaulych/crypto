package args

import (
	"flag"
	"fmt"
)

type RecvArgs struct {
	BindAddr string
}

func ParseRecvCmdFlags(flags *flag.FlagSet, args []string) (*RecvArgs, error) {
	bindHostPtr := flags.String("host", "localhost", "host to bind")
	bindPortPtr := flags.Int("port", 0, "port to bind")

	err := flags.Parse(args)
	if err != nil {
		return nil, err
	}

	addr := fmt.Sprintf("%s:%v", *bindHostPtr, *bindPortPtr)
	return &RecvArgs{BindAddr: addr}, nil
}
