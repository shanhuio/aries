package httpstest

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"

	"shanhu.io/aries/https"
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
	hosts := []string{"127.0.0.1", "::1"}
	hosts = append(hosts, domains...)
	c := &https.RSACertConfig{
		Hosts: hosts,
		IsCA:  true,
	}
	cert, err := https.MakeRSACert(c)
	if err != nil {
		return nil, fmt.Errorf("make RSA cert: %s", err)
	}

	tlsCert, err := cert.X509KeyPair()
	if err != nil {
		return nil, fmt.Errorf("unmarshal TLS cert: %s", err)
	}

	tlsConfig := &tls.Config{
		NextProtos:   []string{"http/1.1"},
		Certificates: []tls.Certificate{tlsCert},
	}

	server := httptest.NewUnstartedServer(h)
	server.TLS = tlsConfig
	server.StartTLS()

	serverURL, err := url.Parse(server.URL)
	if err != nil {
		return nil, fmt.Errorf("invalid server URL: %s", err)
	}

	x509Cert, err := x509.ParseCertificate(tlsCert.Certificate[0])
	if err != nil {
		return nil, fmt.Errorf("parse x509 cert error: %s", err)
	}

	certPool := x509.NewCertPool()
	certPool.AddCert(x509Cert)
	serverHost := serverURL.Host
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{RootCAs: certPool},
		DialContext:     sink(serverHost),
	}

	return &Server{
		Host:      serverHost,
		Server:    server,
		Transport: tr,
	}, nil
}
