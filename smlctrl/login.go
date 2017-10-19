package smlctrl

import (
	"log"

	"shanhu.io/aries/creds"
)

func login(host string, args []string) error {
	server := parseServer(host)
	flags := newFlags()
	flags.Parse(args)

	if _, err := creds.LoginServer(server); err != nil {
		return err
	}

	log.Println("login sucessful, token saved")
	return nil
}
