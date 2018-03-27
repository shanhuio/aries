package ariestest

import (
	"shanhu.io/aries/creds"
	"smallrepo.com/base/httputil"
)

// Login login a server and fetch the token for the given user.
func Login(c *httputil.Client, user, key string) error {
	endPoint := &creds.EndPoint{
		User:        user,
		Server:      c.Server,
		PemFile:     key,
		Transport:   c.Transport,
		Homeless:    true,
		NoTTY:       true,
		NoPermCheck: true,
	}

	login := creds.NewLogin(endPoint)
	token, err := login.Token()
	if err != nil {
		return err
	}

	c.Token = token
	return nil
}
