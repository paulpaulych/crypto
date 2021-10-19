package send

import (
	"github.com/paulpaulych/crypto/cmd/cli"
	"github.com/paulpaulych/crypto/internal/app/messaging/msg-core"
	"github.com/paulpaulych/crypto/internal/app/messaging/protocols"
	"io"
	"math/big"
)

type Conf struct{}

func (conf *Conf) CmdName() string {
	return "send"
}

func (conf *Conf) NewCmd(args []string) (cli.Cmd, cli.CmdConfError) {
	var opts struct {
		P     *cli.BigIntOpt `short:"P" description:"large prime number" required:"true"`
		Input string         `short:"i" long:"input" choice:"file" choice:"console" description:"input type: console or file" default:"console"`
		Args  struct {
			Addr string `positional-arg-name:"target" description:"target host:port. Example: 127.0.0.1:1234"`
			Msg  string `positional-arg-name:"message" description:"message text or name of file when -i=file specified"`
		} `positional-args:"true" required:"true"`
	}

	_, err := cli.ParseFlagsOfCmd(conf.CmdName(), &opts, args)
	if err != nil {
		return nil, err
	}

	msgReader, e := cli.NewInputReader(opts.Input, opts.Args.Msg)
	if e != nil {
		return nil, cli.NewCmdConfErr(e, nil)
	}

	return &Cmd{
		addr: opts.Args.Addr,
		p:    opts.P.Value,
		msg:  msgReader,
	}, nil
}

type Cmd struct {
	addr string
	p    *big.Int
	msg  io.Reader
}

func (cmd *Cmd) Run() error {
	return msg_core.SendMsg(cmd.addr, cmd.msg, protocols.ShamirWriter(cmd.p))
}
