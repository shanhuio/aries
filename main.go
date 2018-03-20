package aries

import (
	"bytes"
	"encoding/json"
	"flag"
	"io"
	"io/ioutil"
	"log"
)

// Main wraps the commin main function for web serving.
type Main struct {
	Addr   string
	Config interface{}
	Serve  func(addr string) error
	Log    *log.Logger
}

// Run runs the service on the given address and config.
func (m *Main) Run(addr string, config io.Reader) error {
	if m.Log == nil {
		m.Log = Log
	}

	if addr == "" {
		addr = m.Addr
	}
	if config != nil {
		dec := json.NewDecoder(config)
		if err := dec.Decode(m.Config); err != nil {
			return err
		}
	}

	return m.Serve(addr)
}

// Main runs the main function body.
func (m *Main) Main() {
	if m.Log == nil {
		m.Log = Log
	}

	addr := flag.String("addr", m.Addr, "address to listen on")
	conf := flag.String("config", "config.json", "config file")
	flag.Parse()

	bs, err := ioutil.ReadFile(*conf)
	if err != nil {
		m.Log.Fatal(err)
	}

	if err := m.Run(*addr, bytes.NewReader(bs)); err != nil {
		m.Log.Fatal(err)
	}
}
