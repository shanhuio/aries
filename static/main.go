package static

import (
	"net/http"

	"shanhu.io/aries"
)

// Config contains the config file for the smlstatic binary.
type Config struct {
	Dir string // Home directory
}

func serve(m *aries.Main) error {
	c := m.Config.(*Config)
	h := http.FileServer(http.Dir(c.Dir))
	s := &http.Server{Handler: h}
	return s.Serve(m.Listener)
}

// Main is the main entrance for smlstatic binary
func Main() {
	c := new(Config)
	m := &aries.Main{
		Addr:   "localhost:8000",
		Config: c,
	}
	m.Main(serve)
}
