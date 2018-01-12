package httpstest

import (
	"net/http"
	"net/http/httptest"
)

// Server wraps a *httptest.Server with HTTP support.
type Server struct {
	*httptest.Server

	Host      string // test host
	Transport *http.Transport
}

// Client creates an HTTP client which transport connects directly to the
// server.
func (s *Server) Client() *http.Client {
	return &http.Client{Transport: s.Transport}
}

// NewServer creates an HTTPS server at the given testing domains.
func NewServer(domains []string, h http.Handler) (*Server, error) {
	c, err := NewTLSConfigs(domains)
	if err != nil {
		return nil, err
	}

	server := httptest.NewUnstartedServer(h)
	server.TLS = c.Server
	server.StartTLS()

	serverHost := server.Listener.Addr().String()
	return &Server{
		Host:      serverHost,
		Server:    server,
		Transport: c.SinkTransport(serverHost),
	}, nil
}
