package recv

import (
	"net"
	"github.com/paulpaulych/crypto/internal/app"
)

type Cmd struct {
	Host   string `short:"h" long:"host" description:"host to bind" default:"localhost"`
	Port   string `short:"p" long:"port" description:"port to bind" default:"12345"`
}

func (c *Cmd) Execute(_ []string) error {
	addr := net.JoinHostPort(c.Host, c.Port)
	return app.ShamirRecv(addr)	
}
