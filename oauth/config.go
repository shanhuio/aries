package oauth

import (
	"time"

	"shanhu.io/aries"
	"shanhu.io/misc/errcode"
)

// JSONConfig is a JSON marshallable config that is commonly used for
// setting up a server.
type JSONConfig struct {
	GitHub       *GitHubApp
	Google       *GoogleApp
	DigitalOcean *DigitalOceanApp
	StateKey     string
	SessionKey   string
	SignInBypass string
	PublicKeys   map[string]string
}

// Config converts a JSON marshallable config to Config.
func (c *JSONConfig) Config() *Config {
	return &Config{
		GitHub:       c.GitHub,
		Google:       c.Google,
		DigitalOcean: c.DigitalOcean,
		StateKey:     []byte(c.StateKey),
		SessionKey:   []byte(c.SessionKey),
		Bypass:       c.SignInBypass,
		KeyStore:     NewFileKeyStore(c.PublicKeys),
	}
}

// SimpleGitHubConfig converts a JSON marshallable config to Config that uses
// Github as the direct user ID mapping. Users that has a public key assigned
// in c.PublicKeys are defined as admin.
func (c *JSONConfig) SimpleGitHubConfig() *Config {
	ret := c.Config()
	ret.SignInCheck = MapGitHub
	ret.Check = func(name string) (interface{}, int, error) {
		if _, isAdmin := c.PublicKeys[name]; isAdmin {
			return nil, 1, nil
		}
		return nil, 0, nil
	}
	return ret
}

// Config is a module configuration for a GitHub Oauth handling module.
type Config struct {
	GitHub       *GitHubApp
	Google       *GoogleApp
	DigitalOcean *DigitalOceanApp

	StateKey        []byte
	SessionKey      []byte
	SessionLifeTime time.Duration
	SessionRefresh  time.Duration

	Bypass   string
	Redirect string

	KeyStore KeyStore

	// SignInCheck exchanges OAuth2 ID's for user ID.
	SignInCheck func(c *aries.C, u *UserMeta, purpose string) (string, error)

	// Check checks the user id and returns the user account structure.
	Check func(user string) (interface{}, int, error)
}

// MapGitHub is a login check function that only allows
// github login. It maps the user ID directly from GitHub users.
func MapGitHub(c *aries.C, u *UserMeta, _ string) (string, error) {
	if u.Method != MethodGitHub {
		return "", errcode.InvalidArgf(
			"login with %q not supported", u.Method,
		)
	}
	return u.ID, nil
}
