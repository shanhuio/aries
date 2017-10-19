package creds

import (
	"os/user"
)

// LoginConfig contains the login stub configuration.
type LoginConfig struct {
	Server  string // Server prefix URL.
	User    string // Optional user name. If blank will fill with OS user name.
	PemFile string // Optional private key. If blank, will use the default key.

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
