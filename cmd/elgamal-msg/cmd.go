package elgamal_msg

import (
	"github.com/paulpaulych/crypto/cmd/elgamal-msg/recv"
	"github.com/paulpaulych/crypto/cmd/elgamal-msg/send"
)

type Cmd struct {
	Send *send.Cmd `command:"send"`
	Recv *recv.Cmd `command:"recv"`
}
