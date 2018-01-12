package httpstest

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"

	"shanhu.io/aries/https"
)

// TLSConfigs creates the certificate setup required for a set of domains.
type TLSConfigs struct {
	Domains []string
	Server  *tls.Config
	Client  *tls.Config
}

// NewTLSConfigs creates a new setup with proper TLS config and HTTP
func NewTLSConfigs(domains []string) (*TLSConfigs, error) {
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

	serverConfig := &tls.Config{
		NextProtos:   []string{"http/1.1"},
		Certificates: []tls.Certificate{tlsCert},
	}

	x509Cert, err := x509.ParseCertificate(tlsCert.Certificate[0])
	if err != nil {
		return nil, fmt.Errorf("parse x509 cert error: %s", err)
	}

	certPool := x509.NewCertPool()
	certPool.AddCert(x509Cert)

	return &TLSConfigs{
		Domains: domains,
		Server:  serverConfig,
		Client:  &tls.Config{RootCAs: certPool},
	}, nil
}

// Sink returns a transport where every outgoing connection dials the same
// sinkAddr but assumes the address is certified as the domains in the
// TLSConfigs.
func (c *TLSConfigs) Sink(sinkAddr string) *http.Transport {
	return &http.Transport{
		DialContext:     sink(sinkAddr),
		TLSClientConfig: c.Client,
	}
}

// SinkHTTPS returns a transport where outgoing https connections dial httpsAddr
// and all other outgoing connections dial httpAddr. When it is HTTPS, it
// assumes the address is certified as the domains in TLSConfigs.
func (c *TLSConfigs) SinkHTTPS(httpAddr, httpsAddr string) *http.Transport {
	return &http.Transport{
		DialContext:     sinkHTTPS(httpAddr, httpsAddr),
		TLSClientConfig: c.Client,
	}
}
