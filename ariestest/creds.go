package ariestest

import (
	"shanhu.io/aries/creds"
	"shanhu.io/misc/httputil"
)

// Login log into a server and fetch the token for the given user.
func Login(c *httputil.Client, user, key string) error {
	endPoint := &creds.Endpoint{
		User:        user,
		Server:      c.Server.String(),
		PemFile:     key,
		Transport:   c.Transport,
		Homeless:    true,
		NoTTY:       true,
		NoPermCheck: true,
	}

	login, err := creds.NewLogin(endPoint)
	if err != nil {
		return err
	}
	token, err := login.Token()
	if err != nil {
		return err
	}

	c.Token = token
	return nil
}
