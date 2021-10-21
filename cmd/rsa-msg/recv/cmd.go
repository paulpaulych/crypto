package recv

import (
	"github.com/paulpaulych/crypto/cmd/cli"
	"github.com/paulpaulych/crypto/internal/app/messaging/msg-core"
	"github.com/paulpaulych/crypto/internal/app/messaging/protocols"
	"github.com/paulpaulych/crypto/internal/app/tcp"
	"net"
)

type Cmd struct {
	Host   string         `short:"h" long:"host" description:"host to bind" default:"localhost"`
	Port   string         `short:"p" long:"port" description:"port to bind" default:"12345"`
	P      *cli.BigIntOpt `short:"P" description:"large prime number"`
	Q      *cli.BigIntOpt `short:"Q" description:"large prime number"`
	Output string         `short:"o" long:"output" choice:"file" choice:"console" description:"output type: console or file" default:"console"`
}

func (c *Cmd) Execute(_ []string) error {
	output, e := cli.NewOutputFactory(c.Output)
	if e != nil {
		return e
	}
	return tcp.StartServer(net.JoinHostPort(c.Host, c.Port),
		msg_core.RecvMessage(
			protocols.RsaReader(c.P.Value, c.Q.Value, output),
		),
	)
}
