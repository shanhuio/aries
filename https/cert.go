package https

import (
	"crypto/tls"
)

// Cert contains a certificate in memory.
type Cert struct {
	Cert []byte // Marshalled PEM block for the certificate.
	Key  []byte // Marshalled PEM block for the private key.
}

// X509KeyPair converts the PEM blocks into a X509 key pair
// for use in an HTTPS server.
func (c *Cert) X509KeyPair() (tls.Certificate, error) {
	return tls.X509KeyPair(c.Cert, c.Key)
}
