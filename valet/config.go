package valet

import (
	"shanhu.io/misc/jsonfile"
)

// Config is the config mapping of the proxy.
type Config struct {
	Admin      string
	PublicKey  string
	SessionKey string

	Cache   string
	Control string
	Hosts   map[string]string
}

func loadConfig(f string) (*Config, error) {
	ret := new(Config)
	if err := jsonfile.Read(f, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

func (c *Config) cache() string {
	if c.Cache == "" {
		return "cache"
	}
	return c.Cache
}

func (c *Config) hosts() []string {
	var ret []string
	for host := range c.Hosts {
		ret = append(ret, host)
	}
	return ret
}

func (c *Config) hasHost(host string) bool {
	if c.Control != "" && host == c.Control {
		return true
	}

	_, ok := c.Hosts[host]
	return ok
}
