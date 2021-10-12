package msg

import (
	. "github.com/paulpaulych/crypto/shamir/cli/cmd/common"
	"github.com/paulpaulych/crypto/shamir/cli/cmd/msg/recv"
	"github.com/paulpaulych/crypto/shamir/cli/cmd/msg/send"
)

type MsgConf struct{}

func (conf *MsgConf) CmdName() string {
	return "msg"
}

func (conf *MsgConf) InitCmd(args []string) (Cmd, CmdConfError) {
	subConfigs := []CmdConf{
		&send.SendConf{},
		&recv.RecvConf{},
	}
	return InitSubCmd(subConfigs, args)
}
