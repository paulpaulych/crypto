package recv

import (
	"net"

	"github.com/paulpaulych/crypto/cmd/cli"
	"github.com/paulpaulych/crypto/internal/app"
)

type Cmd struct {
	Host   string         `short:"h" long:"host" description:"host to bind" default:"localhost"`
	Port   string         `short:"p" long:"port" description:"port to bind" default:"12345"`
	P      *cli.BigIntOpt `short:"P" description:"large prime number" required:"true"`
	G      *cli.BigIntOpt `short:"G" description:"generator of multiplicative group of integers modulo P" required:"true"`
}

func (c *Cmd) Execute(_ []string) error {
	addr := net.JoinHostPort(c.Host, c.Port)
	return app.ElgamalRecv(addr, c.P.Value, c.G.Value)
}
