package rsa_msg

import (
	"github.com/paulpaulych/crypto/cmd/rsa-msg/recv"
	"github.com/paulpaulych/crypto/cmd/rsa-msg/send"
)

type Cmd struct {
	Send *send.Cmd `command:"send"`
	Recv *recv.Cmd `command:"recv"`
}
