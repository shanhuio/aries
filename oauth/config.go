package oauth

import (
	"time"

	"shanhu.io/aries"
	"shanhu.io/misc/errcode"
)

// JSONConfig is a JSON marshallable config that is commonly used for
// setting up a server
type JSONConfig struct {
	GitHub       *GitHubApp
	Google       *GoogleApp
	StateKey     string
	SessionKey   string
	SignInByPass string
	PublicKeys   map[string]string
}

// Config converts a JSON marshallable config to Config.
func (c *JSONConfig) Config() *Config {
	return &Config{
		GitHub:     c.GitHub,
		Google:     c.Google,
		StateKey:   []byte(c.StateKey),
		SessionKey: []byte(c.SessionKey),
		ByPass:     c.SignInByPass,
		KeyStore:   NewFileKeyStore(c.PublicKeys),
	}
}

// GitHubBasedConfig converts a JSON marshallable config to Config that uses
// Github as the direct user ID mapping. Users that has a public key assigned
// in c.PublicKeys are defined as admin.
func (c *JSONConfig) GitHubBasedConfig() *Config {
	ret := c.Config()
	ret.LoginCheck = MapGitHub
	return ret
}

// Config is a module configuration for a GitHub Oauth handling module.
type Config struct {
	GitHub *GitHubApp
	Google *GoogleApp

	StateKey        []byte
	SessionKey      []byte
	SessionLifeTime time.Duration
	SessionRefresh  time.Duration

	ByPass   string
	Redirect string

	KeyStore KeyStore

	// Exchanges Oauth2 ID's for user ID.
	LoginCheck func(c *aries.C, method, id string) (string, error)

	// Fetches the user account structure.
	Check func(user string) (interface{}, int)
}

// MapGitHub is a login check function that only allows
// github login. It maps the user ID directly from GitHub users.
func MapGitHub(c *aries.C, method, id string) (string, error) {
	if method != "github" {
		return "", errcode.InvalidArgf(
			"login with %q not supported", method,
		)
	}
	return id, nil
}
