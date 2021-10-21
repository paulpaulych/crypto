package shamir_msg

import (
	"github.com/paulpaulych/crypto/cmd/shamir-msg/recv"
	"github.com/paulpaulych/crypto/cmd/shamir-msg/send"
)

type Cmd struct {
	Send *send.Cmd `command:"send"`
	Recv *recv.Cmd `command:"recv"`
}
