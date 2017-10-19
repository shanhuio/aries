package aries

import (
	"fmt"
	"strings"

	"shanhu.io/misc/signer"
)

// Bearer returns the authorization token.
func Bearer(c *C) string {
	auth := c.Req.Header.Get("Authorization")
	if !strings.HasPrefix(auth, "Bearer ") {
		return ""
	}
	return strings.TrimPrefix(auth, "Bearer ")
}

// CheckToken checks if the bearer token is properly signed by the
// same API key.
func CheckToken(c *C, s *signer.TimeSigner) error {
	token := Bearer(c)
	if !s.Check(token) {
		return fmt.Errorf("invalid token")
	}

	return nil
}
