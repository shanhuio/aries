package oauth

import (
	"time"

	"shanhu.io/aries"
)

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
	LoginCheck func(c *aries.C, method, id string) string

	// Fetches the user account structure.
	Check func(user string) (interface{}, bool)
}
