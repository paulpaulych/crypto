package send

import (
	"github.com/paulpaulych/crypto/cmd/cli"
	"github.com/paulpaulych/crypto/internal/app/messaging/msg-core"
	"github.com/paulpaulych/crypto/internal/app/messaging/protocols"
)

type Cmd struct {
	BobPub string `short:"b" long:"bob-pub" description:"path to file containing destination public key" required:"true"`
	Input  string `short:"i" long:"input" choice:"file" choice:"console" description:"input type: console or file" default:"console"`
	Args   struct {
		Addr string `positional-arg-name:"target address" description:"target host:port. Example: 127.0.0.1:1234"`
		Msg  string `positional-arg-name:"message" description:"message text or name of file when -i=file specified"`
	} `positional-args:"true" required:"true"`
}

func (c *Cmd) Execute(_ []string) error {
	msgReader, e := cli.NewInputReader(c.Input, c.Args.Msg)
	if e != nil {
		return e
	}
	return msg_core.SendMsg(c.Args.Addr, msgReader, protocols.RsaWriter(c.BobPub))
}
