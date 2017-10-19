package oauth

import (
	"crypto/rsa"
	"fmt"

	"golang.org/x/crypto/ssh"

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

func unmarshalPublicKey(bs []byte) (*rsa.PublicKey, error) {
	if len(bs) == 0 {
		return nil, fmt.Errorf("public key not present")
	}
	k, _, _, _, err := ssh.ParseAuthorizedKey(bs)
	if err != nil {
		return nil, err
	}
	if k.Type() != "ssh-rsa" {
		return nil, fmt.Errorf("public key is not an RSA key")
	}
	ck, ok := k.(ssh.CryptoPublicKey)
	if !ok {
		return nil, fmt.Errorf("public key is not an RSA key")
	}

	return ck.CryptoPublicKey().(*rsa.PublicKey), nil
}
