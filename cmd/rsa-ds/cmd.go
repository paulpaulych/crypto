package rsa_ds

import (
	"github.com/paulpaulych/crypto/cmd/rsa-ds/key-gen"
	"github.com/paulpaulych/crypto/cmd/rsa-ds/sign"
	"github.com/paulpaulych/crypto/cmd/rsa-ds/validate"
)

type Cmd struct {
	KeyGen   *key_gen.Cmd  `command:"key-gen"`
	Sign     *sign.Cmd     `command:"sign"`
	Validate *validate.Cmd `command:"validate"`
}
