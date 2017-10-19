package smlctrl

import (
	"shanhu.io/misc/subcmd"
)

func cmds() *subcmd.List {
	ret := subcmd.New()
	ret.AddHost("login", "login and fetch the token mint", login)
	ret.AddHost("deploy", "deploy a service instance", deploy)
	return ret
}

// Main is the main entrance for smlctrl utility.
func Main() { cmds().Main() }
