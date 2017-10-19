package creds

import (
	"shanhu.io/aries/oauth"
)

// Creds is the credential that is cached after logging in.
type Creds struct {
	Server      string
	oauth.Creds // user is saved in creds.
}
