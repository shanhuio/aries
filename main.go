package aries

import (
	"flag"

	"shanhu.io/misc/jsonfile"
)

// Main wraps the commin main function for web serving.
type Main struct {
	Addr   string
	Config interface{}
	Logger *Logger
}

// Main runs the main function body.
func (m *Main) Main(serve func(m *Main) error) {
	if m.Logger == nil {
		m.Logger = StdLogger()
	}

	flag.StringVar(&m.Addr, "addr", m.Addr, "address to listen on")
	conf := flag.String("config", "config.json", "config file")
	flag.Parse()

	if err := jsonfile.Read(*conf, m.Config); err != nil {
		m.Logger.Exit(err)
	}

	if err := serve(m); err != nil {
		m.Logger.Exit(err)
	}
}
