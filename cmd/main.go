package main

import (
	"github.com/jessevdk/go-flags"
	"github.com/paulpaulych/crypto/cmd/elgamal-msg"
	"github.com/paulpaulych/crypto/cmd/rsa-ds"
	"github.com/paulpaulych/crypto/cmd/rsa-msg"
	"github.com/paulpaulych/crypto/cmd/shamir-msg"
)

func main() {
	var rootCmd struct {
		RsaMsg     *rsa_msg.Cmd     `command:"rsa-msg" description:"allows messaging encrypted by RSA algorithm"`
		ElgamalMsg *elgamal_msg.Cmd `command:"elgamal-msg" description:"allows messaging encrypted by Elgamal algorithm"`
		ShamirMsg  *shamir_msg.Cmd  `command:"shamir-msg" description:"allows messaging encrypted by Shamir algorithm"`
		RsaDs      *rsa_ds.Cmd      `command:"rsa-ds" description:"for work with RSA based digital signature"`
	}
	parser := flags.NewParser(&rootCmd, flags.Default)
	_, err := parser.Parse()
	if err != nil {
		return
	}
}
