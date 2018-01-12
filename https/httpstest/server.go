package httpstest

import (
	"net/http"
	"net/http/httptest"
)

// Server wraps a *httptest.Server with HTTP support.
type Server struct {
	*httptest.Server

	Transport *http.Transport
}

// Client creates an HTTP client which transport connects directly to the
// server.
func (s *Server) Client() *http.Client {
	return &http.Client{Transport: s.Transport}
}

// NewServer creates an HTTPS server for the given testing domains.
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
		Server:    server,
		Transport: c.Sink(serverHost),
	}, nil
}

// DualServer wraps two *httptest.Server's with a transport that
// goes to one of them base on HTTP or HTTPS.
type DualServer struct {
	HTTP      *httptest.Server
	HTTPS     *httptest.Server
	Transport *http.Transport
}

// NewDualServer creates an HTTPS dual server for the given testing domains.
func NewDualServer(domains []string, h http.Handler) (*DualServer, error) {
	c, err := NewTLSConfigs(domains)
	if err != nil {
		return nil, err
	}

	httpsServer := httptest.NewUnstartedServer(h)
	httpsServer.TLS = c.Server
	httpsServer.StartTLS()

	httpServer := httptest.NewServer(h)

	httpAddr := httpServer.Listener.Addr().String()
	httpsAddr := httpsServer.Listener.Addr().String()
	return &DualServer{
		HTTP:      httpServer,
		HTTPS:     httpsServer,
		Transport: c.SinkHTTPS(httpAddr, httpsAddr),
	}, nil
}

// Client creates an HTTP client which transport connects directly to one
// of the servers base on the protocol port.
func (s *DualServer) Client() *http.Client {
	return &http.Client{Transport: s.Transport}
}
