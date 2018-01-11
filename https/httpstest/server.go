package httpstest

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
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

	serverURL, err := url.Parse(server.URL)
	if err != nil {
		return nil, fmt.Errorf("invalid server URL: %s", err)
	}
	serverHost := serverURL.Host

	return &Server{
		Host:      serverHost,
		Server:    server,
		Transport: c.SinkTransport(serverHost),
	}, nil
}
