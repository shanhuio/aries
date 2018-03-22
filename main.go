package aries

import (
	"flag"
	"net"

	"shanhu.io/misc/jsonfile"
)

// Main wraps the commin main function for web serving.
type Main struct {
	Addr     string // default address to listen
	Listener net.Listener
	Config   interface{}
	Logger   *Logger
}

// Main runs the main function body.
func (m *Main) Main(serve func(m *Main) error) {
	if m.Logger == nil {
		m.Logger = StdLogger()
	}

	if m.Listener == nil {
		flag.StringVar(&m.Addr, "addr", m.Addr, "address to listen on")
	} else {
		m.Addr = "" // will use the given listener
	}
	conf := flag.String("config", "config.json", "config file")
	flag.Parse()

	if err := jsonfile.Read(*conf, m.Config); err != nil {
		m.Logger.Exit(err)
	}

	if m.Listener == nil {
		lis, err := net.Listen("tcp", m.Addr)
		if err != nil {
			m.Logger.Exit(err)
		}
		m.Listener = lis

		m.Logger.Printf("serving on %s", m.Addr)
	}

	if err := serve(m); err != nil {
		m.Logger.Exit(err)
	}
}
