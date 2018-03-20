package aries

import (
	"bytes"
	"encoding/json"
	"flag"
	"io"
	"io/ioutil"
)

// Main wraps the commin main function for web serving.
type Main struct {
	Addr   string
	Config interface{}
	Serve  func(addr string, logger *Logger) error
	Logger *Logger
}

// Run runs the service on the given address and config.
func (m *Main) Run(config io.Reader) error {
	if config != nil {
		dec := json.NewDecoder(config)
		if err := dec.Decode(m.Config); err != nil {
			return err
		}
	}

	return m.Serve(m.Addr, m.Logger)
}

// Main runs the main function body.
func (m *Main) Main() {
	if m.Logger == nil {
		m.Logger = StdLogger()
	}

	flag.StringVar(&m.Addr, "addr", m.Addr, "address to listen on")
	conf := flag.String("config", "config.json", "config file")
	flag.Parse()

	bs, err := ioutil.ReadFile(*conf)
	if err != nil {
		m.Logger.Exit(err)
	}

	if err := m.Run(bytes.NewReader(bs)); err != nil {
		m.Logger.Exit(err)
	}
}
