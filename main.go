package aries

import (
	"context"
	"flag"
	"log"
	"net/http"
	"strings"

	"shanhu.io/misc/jsonutil"
	"shanhu.io/misc/jsonx"
	"shanhu.io/misc/unixhttp"
)

func runMain(
	b BuildFunc, configFile string, config interface{}, addr string,
) error {
	if config != nil {
		if strings.HasSuffix(configFile, ".json") {
			if err := jsonutil.ReadFile(configFile, config); err != nil {
				return err
			}
		} else {
			if err := jsonx.ReadFile(configFile, config); err != nil {
				return err
			}
		}
	}

	s, err := b(&Env{
		Context: context.Background(),
		Config:  config,
	})
	if err != nil {
		return err
	}

	log.Printf("serve on %s", addr)

	if strings.HasSuffix(addr, ".sock") {
		return unixhttp.ListenAndServe(addr, Serve(s))
	}
	return http.ListenAndServe(addr, Serve(s))
}

// Main launches a service with the given config structure, and default
// address.
func Main(b BuildFunc, config interface{}, addr string) {
	flag.StringVar(&addr, "addr", addr, "address to listen on")
	var configFile string
	if config != nil {
		flag.StringVar(&configFile, "config", "config.json", "config file")
	}
	flag.Parse()

	if err := runMain(b, configFile, config, addr); err != nil {
		log.Fatal(err)
	}
}

// SimpleMain launches a service with no config and default address.
func SimpleMain(service Service, addr string) {
	f := func(_ *Env) (Service, error) { return service, nil }
	Main(f, nil, addr)
}
