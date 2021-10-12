package send

import (
	"flag"
	"fmt"
	. "github.com/paulpaulych/crypto/shamir/cli/cmd/common"
	"github.com/paulpaulych/crypto/shamir/cli/internal/messaging"
	"github.com/paulpaulych/crypto/shamir/cli/internal/messaging/protocols"
	"math/big"
)

type SendConf struct{}

func (conf *SendConf) CmdName() string {
	return "send"
}

func (conf *SendConf) InitCmd(args []string) (Cmd, CmdConfError) {
	flags := flag.NewFlagSet(conf.CmdName(), flag.ContinueOnError)

	protocolPtr := flags.String("protocol", "", "protocol")
	primePtr := flags.String("prime", "", "prime integer")

	err := Parse(flags, args)
	if err != nil {
		return nil, err
	}

	if flags.NArg() < 2 {
		return nil, NewCmdConfError("args required: [host:port] [message]", nil)
	}

	addr, msgStr := flags.Arg(0), flags.Arg(1)
	msg, success := new(big.Int).SetString(msgStr, 10)
	if !success {
		return nil, NewCmdConfError("message must be integer", nil)
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

func writerForProtocol(name string, primeStr string) (messaging.WriteMsg, CmdConfError) {
	switch name {
	case "shamir":
		if len(primeStr) == 0 {
			return nil, NewCmdConfError("shamir protocol requires -prime flag", nil)
		}
		prime, success := new(big.Int).SetString(primeStr, 10)
		if !success {
			return nil, NewCmdConfError("cannot parse prime", nil)
		}
		return protocols.ShamirWriter(prime), nil
	default:
		msg := fmt.Sprintf("unknown protocol '%s'", name)
		return nil, NewCmdConfError(msg, nil)
	}
}
