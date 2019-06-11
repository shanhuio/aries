package aries

import (
	"flag"
	"net"
	"net/http"
	"strings"

	"shanhu.io/misc/jsonutil"
)

// Main launches a service with the given config structure, and default
// address.
func Main(b BuildFunc, config interface{}, addr string) {
	flag.StringVar(&addr, "addr", addr, "address to listen on")
	conf := flag.String("config", "config.json", "config file")
	flag.Parse()

	logger := StdLogger()
	if err := jsonutil.ReadFile(*conf, config); err != nil {
		logger.Exit(err)
	}

	s, err := b(&Env{
		Config: config,
		Logger: logger,
	})
	if err != nil {
		logger.Exit(err)
	}

	logger.Printf("serve on %s", addr)

	if strings.HasSuffix(addr, ".sock") {
		lis, err := net.ListenUnix("unix", &net.UnixAddr{
			Name: addr,
			Net:  "unix",
		})
		if err != nil {
			logger.Exit(err)
		}
		logger.Exit(http.Serve(lis, Serve(s)))
	}

	logger.Exit(http.ListenAndServe(addr, Serve(s)))
}
