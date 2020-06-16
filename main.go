package aries

import (
	"context"
	"flag"
	"log"
	"net/http"
	"strings"

	"shanhu.io/misc/errcode"
	"shanhu.io/misc/jsonutil"
	"shanhu.io/misc/jsonx"
	"shanhu.io/misc/osutil"
	"shanhu.io/misc/unixhttp"
)

func loadConfig(file string, config interface{}) error {
	if file == "" {
		for _, try := range []string{
			"config.jsonx",
			"config.json",
		} {
			ok, err := osutil.IsRegular(try)
			if err != nil {
				return err
			}
			if ok {
				file = try
				break
			}
		}
	}
	if file == "" {
		return errcode.InvalidArgf("config file not specified")
	}
	log.Println("config filename: ", file)

	if strings.HasSuffix(file, ".json") {
		return jsonutil.ReadFile(file, config)
	}
	return jsonx.ReadFile(file, config)
}

func runMain(
	b BuildFunc, configFile string, config interface{}, addr string,
) error {
	if config != nil {
		if err := loadConfig(configFile, config); err != nil {
			return errcode.Annotate(err, "load config file")
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
		flag.StringVar(&configFile, "config", "", "config file")
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
