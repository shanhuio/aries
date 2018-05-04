package oauth

import (
	"shanhu.io/misc/signer"
)

// LoginRequest is the request for logging in.
type LoginRequest struct {
	User       string
	SignedTime *signer.SignedRSABlock
}

// Creds is the response for logging in.
type Creds struct {
	User    string
	Token   string
	Expires int64 // nanosecond timestamp
}
