package rsa_ds

import (
	"github.com/paulpaulych/crypto/cmd/cli"
	key_gen "github.com/paulpaulych/crypto/cmd/rsa-ds/key-gen"
	"github.com/paulpaulych/crypto/cmd/rsa-ds/sign"
	"github.com/paulpaulych/crypto/cmd/rsa-ds/validate"
)

type Conf struct{}

func (conf *Conf) CmdName() string {
	return "rsa-ds"
}

func (conf *Conf) NewCmd(args []string) (cli.Cmd, cli.CmdConfError) {
	subConfigs := []cli.CmdConf{
		&key_gen.Conf{},
		&sign.Conf{},
		&validate.Conf{},
	}
	return cli.InitSubCmd(subConfigs, args)
}
