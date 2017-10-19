package smlctrl

import (
	"fmt"
	"log"
	"os"

	"shanhu.io/aries/creds"
)

func deploy(host string, args []string) error {
	server := parseServer(host)
	flags := newFlags()
	flags.Parse(args)
	args = flags.Args()

	if len(args) == 0 {
		return fmt.Errorf("require instance name")
	}

	c, err := creds.Dial(server)
	if err != nil {
		return err
	}

	for _, arg := range args {
		log.Println(arg)

		req := struct {
			Name string
		}{
			Name: arg,
		}

		if err := c.JSONPost("/deploy", &req, os.Stdout); err != nil {
			return err
		}
	}
	return nil
}
