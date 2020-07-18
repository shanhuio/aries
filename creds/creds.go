package creds

import (
	"crypto/rsa"
	"fmt"
	"net/http"

	"shanhu.io/aries/oauth"
	"shanhu.io/misc/httputil"
	"shanhu.io/misc/signer"
)

// Creds is the credential that is cached after logging in. This can also be
// saved in JSON format in user's home directory.
type Creds struct {
	Server      string
	oauth.Creds // user is saved in creds.
}

// Request contains the configuration to create a credential.
type Request struct {
	Server string
	User   string
	Key    *rsa.PrivateKey
	TTL    time.Duration

	// Transport is the http transport for the token exchange.
	Transport http.RoundTripper
}

// NewCredsFromRequest creates a new user credential by dialing the server
// using the given RSA private key.
func NewCredsFromRequest(req *Request) (*Creds, error) {
	signed, err := signer.RSASignTime(req.Key)
	if err != nil {
		return nil, err
	}

	login := &oauth.LoginRequest{
		User:       req.User,
		SignedTime: signed,
		TTL:        req.TTL,
	}
	cs := &Creds{Server: req.Server}

	c, err := httputil.NewClient(req.Server)
	if err != nil {
		return nil, err
	}
	c.Transport = req.Transport

	if err := c.JSONCall("/pubkey/signin", login, &cs.Creds); err != nil {
		return nil, err
	}

	if got := cs.Creds.User; got != req.User {
		return nil, fmt.Errorf("login as user %q, got %q", req.User, got)
	}

	return cs, nil
}

// NewCreds creates a new user credential by dialing the server using
// the given RSA private key.
func NewCreds(server, user string, k *rsa.PrivateKey) (*Creds, error) {
	req := &Request{
		Server: server,
		User:   user,
		Key:    k,
	}
	return NewCredsFromRequest(req)
}
