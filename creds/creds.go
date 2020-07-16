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

// NewCreds creates a new credential
func NewCreds(
	server, user string, k *rsa.PrivateKey, tr http.RoundTripper,
) (*Creds, error) {
	signed, err := signer.RSASignTime(k)
	if err != nil {
		return nil, err
	}

	req := &oauth.LoginRequest{
		User:       user,
		SignedTime: signed,
	}
	cs := &Creds{Server: server}

	c, err := httputil.NewClient(server)
	if err != nil {
		return nil, err
	}
	if tr != nil {
		c.Transport = tr
	}
	if err := c.JSONCall("/pubkey/signin", req, &cs.Creds); err != nil {
		return nil, err
	}

	if got := cs.Creds.User; got != user {
		return nil, fmt.Errorf("login as user %q, got %q", user, got)
	}

	return cs, nil
}
