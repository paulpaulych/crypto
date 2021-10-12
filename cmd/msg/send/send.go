package send

import (
	"flag"
	"fmt"
	"github.com/paulpaulych/crypto/internal/app/messaging"
	"github.com/paulpaulych/crypto/internal/app/messaging/protocols"
	cli2 "github.com/paulpaulych/crypto/internal/infra/cli"
	"math/big"
)

type SendConf struct{}

func (conf *SendConf) CmdName() string {
	return "send"
}

func (conf *SendConf) InitCmd(args []string) (cli2.Cmd, cli2.CmdConfError) {
	flags := flag.NewFlagSet(conf.CmdName(), flag.ContinueOnError)

	protocolPtr := flags.String("protocol", "", "protocol")
	primePtr := flags.String("prime", "", "prime integer")

	err := cli2.Parse(flags, args)
	if err != nil {
		return nil, err
	}

	if flags.NArg() < 2 {
		return nil, cli2.NewCmdConfError("args required: [host:port] [message]", nil)
	}

	addr, msgStr := flags.Arg(0), flags.Arg(1)
	msg, success := new(big.Int).SetString(msgStr, 10)
	if !success {
		return nil, cli2.NewCmdConfError("message must be integer", nil)
	}

	writer, err := writerForProtocol(*protocolPtr, *primePtr)
	if err != nil {
		return nil, err
	}
	return &SendCmd{addr: addr, writer: writer, msg: msg}, nil
}

type SendCmd struct {
	addr   string
	writer messaging.WriteMsg
	msg    messaging.Msg
}

func (cmd *SendCmd) Run() error {
	return messaging.SendMsg(cmd.addr, cmd.msg, cmd.writer)
}

func writerForProtocol(name string, primeStr string) (messaging.WriteMsg, cli2.CmdConfError) {
	switch name {
	case "shamir":
		if len(primeStr) == 0 {
			return nil, cli2.NewCmdConfError("shamir protocol requires -prime flag", nil)
		}
		prime, success := new(big.Int).SetString(primeStr, 10)
		if !success {
			return nil, cli2.NewCmdConfError("cannot parse prime", nil)
		}
		return protocols.ShamirWriter(prime), nil
	default:
		msg := fmt.Sprintf("unknown protocol '%s'", name)
		return nil, cli2.NewCmdConfError(msg, nil)
	}
}
