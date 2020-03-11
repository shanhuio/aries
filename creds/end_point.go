package creds

import (
	"net/http"
	"os"
	"os/user"

	"shanhu.io/aries"
)

// EndPoint contains the login stub configuration.
type EndPoint struct {
	// Server prefix URL.
	Server string

	// Optional user name. If blank will fill with OS user name.
	User string

	// Optional private key. If blank, will use the default key.
	PemFile string

	// Optional transport for creating the client.
	Transport http.RoundTripper

	Homeless bool // If true, will not look into the home folder for caches.
	NoTTY    bool // If true, will not fail if the key is encrypted.

	// If skip checking key permission
	NoPermCheck bool
}

func currentUser() (string, error) {
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

// NewEndPoint creates a new default endpoint for the target server.
func NewEndPoint(server string) (*EndPoint, error) {
	user, err := currentUser()
	if err != nil {
		return nil, err
	}
	return &EndPoint{User: user, Server: server}, nil
}

// NewRobot creates a new robot endpoint.
func NewRobot(user, server, key string, env *aries.Env) *EndPoint {
	ep := &EndPoint{
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
