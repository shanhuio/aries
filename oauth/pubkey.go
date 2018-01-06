package oauth

import (
	"crypto/rsa"
	"errors"
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

var errNotRSA = errors.New("public key is not an RSA key")

func unmarshalPublicKey(bs []byte) (*rsa.PublicKey, error) {
	if len(bs) == 0 {
		return nil, fmt.Errorf("public key not present")
	}
	k, _, _, _, err := ssh.ParseAuthorizedKey(bs)
	if err != nil {
		return nil, err
	}

	if k.Type() != "ssh-rsa" {
		return nil, errNotRSA
	}
	ck, ok := k.(ssh.CryptoPublicKey)
	if !ok {
		return nil, errNotRSA
	}

	ret, ok := ck.CryptoPublicKey().(*rsa.PublicKey)
	if !ok {
		return nil, errNotRSA
	}
	return ret, nil
}
