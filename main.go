package aries

import (
	"context"
	"flag"
	"net/http"
	"strings"

	"shanhu.io/misc/jsonutil"
	"shanhu.io/misc/unixhttp"
)

// Main launches a service with the given config structure, and default
// address.
func Main(b BuildFunc, config interface{}, addr string) {
	flag.StringVar(&addr, "addr", addr, "address to listen on")
	var conf string
	if config != nil {
		flag.StringVar(&conf, "config", "config.json", "config file")
	}
	flag.Parse()

	logger := StdLogger()
	if config != nil {
		if err := jsonutil.ReadFile(conf, config); err != nil {
			logger.Exit(err)
		}
	}

	s, err := b(&Env{
		Context: context.Background(),
		Config:  config,
		Logger:  logger,
	})
	if err != nil {
		logger.Exit(err)
	}

	logger.Printf("serve on %s", addr)

	if strings.HasSuffix(addr, ".sock") {
		logger.Exit(unixhttp.ListenAndServe(addr, Serve(s)))
	}

	logger.Exit(http.ListenAndServe(addr, Serve(s)))
}
