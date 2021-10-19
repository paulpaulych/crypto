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
	var opts struct {
		P          *cli.BigIntOpt `short:"P" description:"large prime number" required:"true"`
		G          *cli.BigIntOpt `short:"G" description:"generator of multiplicative group of integers modulo P" required:"true"`
		BobPubFile string         `short:"b" long:"bob-pub" description:"path to file containing destination public key" required:"true"`
		Input      string         `short:"i" long:"input" choice:"file" choice:"console" description:"input type: console or file" default:"console"`
		Args       struct {
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
		addr:       opts.Args.Addr,
		p:          opts.P.Value,
		g:          opts.G.Value,
		bobPubFile: opts.BobPubFile,
		msg:        msgReader,
	}, nil
}

type Cmd struct {
	addr       string
	p, g       *big.Int
	bobPubFile string
	msg        io.Reader
}

func (cmd *Cmd) Run() error {
	commonPub, e := dh.NewCommonPublicKey(cmd.p, cmd.g)
	if e != nil {
		return fmt.Errorf("Diffie-Hellman public key error: %v", e)
	}
	return msg_core.SendMsg(cmd.addr, cmd.msg, protocols.ElgamalWriter(commonPub, cmd.bobPubFile))
}
