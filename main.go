package aries

import (
	"flag"
	"net/http"

	"shanhu.io/misc/jsonfile"
)

// Main launches a service with the given config structure, and default
// address.
func Main(b BuildFunc, config interface{}, addr string) {
	flag.StringVar(&addr, "addr", addr, "address to listen on")
	conf := flag.String("config", "config.json", "config file")
	flag.Parse()

	logger := StdLogger()
	if err := jsonfile.Read(*conf, config); err != nil {
		logger.Exit(err)
	}

	s, err := b(&Env{
		Config: config,
		Logger: logger,
	})
	if err != nil {
		logger.Exit(err)
	}

	logger.Exit(http.ListenAndServe(addr, Serve(s)))
}
