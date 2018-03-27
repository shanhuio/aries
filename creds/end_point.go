package creds

import (
	"os/user"

	"net/http"
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
}

func currentUser() (string, error) {
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
