package smlctrl

import (
	"flag"

	"smallrepo.com/base/httputil"
)

func newFlags() *flag.FlagSet {
	return flag.NewFlagSet("smlctrl", flag.ExitOnError)
}

func parseServer(s string) string {
	if s == "" {
		return "https://ctrl.shanhu.io"
	}
	return httputil.ExtendServer(s)
}
