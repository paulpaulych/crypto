package flags

import (
	"errors"
	"flag"
	"fmt"
	cmd "github.com/paulpaulych/crypto/shamir/cli/cmd"
	"math/big"
)

func ParseSendCmdFlags(flags *flag.FlagSet, args []string) (*cmd.SendArgs, error) {
	targetHostPtr := flags.String("host", "localhost", "target host")
	targetPortPtr := flags.Int("port", 0, "target port")
	largePrimePtr := flags.String("P", "", "prime integer")
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
	msg, success := new(big.Int).SetString(*msgPtr, 10)
	if !success {
		return nil, errors.New("cannot parse message")
	}
	P, success := new(big.Int).SetString(*largePrimePtr, 10)
	if !success {
		return nil, errors.New("cannot parse P")
	}
	return &cmd.SendArgs{Addr: addr, Msg: msg, P: P}, nil
}
