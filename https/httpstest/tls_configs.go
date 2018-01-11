package httpstest

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"

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
