package aries

import (
	"flag"

	"shanhu.io/misc/jsonfile"
)

// Main wraps the commin main function for web serving.
type Main struct {
	Addr   string
	Config interface{}
	Serve  func(addr string) error
}

// Main runs the main function body.
func (m *Main) Main() {
	flag.StringVar(&m.Addr, "addr", m.Addr, "address to listen on")
	conf := flag.String("config", "config.json", "config file")
	flag.Parse()

	if err := jsonfile.Read(*conf, m.Config); err != nil {
		Log.Fatal(err)
	}

	if err := m.Serve(m.Addr); err != nil {
		Log.Fatal(err)
	}
}
