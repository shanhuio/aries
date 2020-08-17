package creds

import (
	"net/http"
	"os"
	"os/user"

	"shanhu.io/aries"
)

// Endpoint contains the login stub configuration.
type Endpoint struct {
	// Server is the server's prefix URL.
	Server string

	// User is an optional user name. If blank will use OS user name, or the
	// value of SHANHU_USER environment variable if exists.
	User string

	// Optional private key content. If nil, will use fall to use
	// PemFile. When presented, PemFile is ignored.
	Key []byte

	// Optional private key. If blank, will use the default key.
	PemFile string

	// Optional transport for creating the client.
	Transport http.RoundTripper

	Homeless bool // If true, will not look into the home folder for caches.
	NoTTY    bool // If true, will not fail if the key is encrypted.

	// If skip checking key permission
	NoPermCheck bool
}

// CurrentUser returns the new name of current user.
func CurrentUser() (string, error) {
	v, ok := os.LookupEnv("SHANHU_USER")
	if ok {
		return v, nil
	}

	u, err := user.Current()
	if err != nil {
		return "", err
	}
	return u.Username, nil
}

// NewEndpoint creates a new default endpoint for the target server.
func NewEndpoint(server string) (*Endpoint, error) {
	user, err := CurrentUser()
	if err != nil {
		return nil, err
	}
	return &Endpoint{User: user, Server: server}, nil
}

// NewRobot creates a new robot endpoint.
func NewRobot(user, server, key string, env *aries.Env) *Endpoint {
	ep := &Endpoint{
		Server:   server,
		User:     user,
		PemFile:  key,
		Homeless: true,
		NoTTY:    true,
	}
	if env != nil {
		ep.Transport = env.Transport
		ep.NoPermCheck = env.Testing
	}
	return ep
}
