package args

import (
	"errors"
	"flag"
	"fmt"
)

type SendArgs struct {
	Addr string
	Msg  string
}

func ParseSendCmdFlags(flags *flag.FlagSet, args []string) (*SendArgs, error) {
	targetHostPtr := flags.String("host", "localhost", "target host")
	targetPortPtr := flags.Int("port", 0, "target port")
	msgPtr := flags.String("msg", "", "Integer message to send")

	err := flags.Parse(args)
	if err != nil {
		flags.PrintDefaults()
		return nil, err
	}
	if *msgPtr == "" {
		return nil, errors.New("msg is empty")
	}
	addr := fmt.Sprintf("%s:%v", *targetHostPtr, *targetPortPtr)
	return &SendArgs{Addr: addr, Msg: *msgPtr}, nil
}
