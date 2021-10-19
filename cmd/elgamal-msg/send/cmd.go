package send

import (
	"fmt"
	"github.com/paulpaulych/crypto/cmd/cli"
	"github.com/paulpaulych/crypto/internal/app/messaging/msg-core"
	"github.com/paulpaulych/crypto/internal/app/messaging/protocols"
	dh "github.com/paulpaulych/crypto/internal/core/diffie-hellman"
	"io"
	"math/big"
)

type Conf struct{}

func (conf *Conf) CmdName() string {
	return "send"
}

func (conf *Conf) NewCmd(args []string) (cli.Cmd, cli.CmdConfError) {
	flagsSpec := cli.NewFlagSpec(conf.CmdName(), map[string]string{
		"P":       "large prime number",
		"G":       "generator of multiplicative group of integers modulo P",
		"bob-pub": "path to file containing destination public key",
		"i":       "message input type: file or arg",
	})

	flags, err := flagsSpec.Parse(args)
	if err != nil {
		return nil, err
	}

	if len(flags.Args) < 2 {
		return nil, cli.NewCmdConfError("args required: [host:port] [message]", nil)
	}

	addr := flags.Args[0]

	bobPubFName := flags.Flags["bob-pub"].Get()
	input := flags.Flags["i"].GetOr("console")

	msgReader, e := cli.NewInputReader(input, flags.Args[1:])
	if e != nil {
		return nil, cli.NewCmdConfError(e.Error(), nil)
	}

	pStr := flags.Flags["P"].Get()
	gStr := flags.Flags["G"].Get()
	P, success := new(big.Int).SetString(*pStr, 10)
	if !success {
		return nil, cli.NewCmdConfError("cannot parse P", nil)
	}
	G, success := new(big.Int).SetString(*gStr, 10)
	if !success {
		return nil, cli.NewCmdConfError("cannot parse G", nil)
	}
	if bobPubFName == nil {
		return nil, cli.NewCmdConfError("required flag: -bob-pub", nil)
	}
	commonPub, e := dh.NewCommonPublicKey(P, G)
	if e != nil {
		return nil, cli.NewCmdConfError(
			fmt.Sprintf("Diffie-Hellman public key error: %v", e), nil,
		)
	}
	return &Cmd{addr: addr, alice: protocols.ElgamalWriter(commonPub, *bobPubFName), msg: msgReader}, nil
}

type Cmd struct {
	addr  string
	alice msg_core.ConnWriter
	msg   io.Reader
}

func (cmd *Cmd) Run() error {
	return msg_core.SendMsg(cmd.addr, cmd.msg, cmd.alice)
}
