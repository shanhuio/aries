package creds

import (
	"shanhu.io/aries/oauth"
)

// Creds is the credential that is cached after logging in. This can also be
// saved in JSON format in user's home directory.
type Creds struct {
	Server      string
	oauth.Creds // user is saved in creds.
}
