package identity

import (
	"net/url"

	"shanhu.io/misc/errcode"
	"shanhu.io/misc/httputil"
	"shanhu.io/misc/rsautil"
)

// WebKeyRegistry is a storage of public keys backed by a web site.
type WebKeyRegistry struct {
	client *httputil.Client
}

// NewWebKeyRegistry creates a new key store backed by a web site
// at the given base URL.
func NewWebKeyRegistry(base *url.URL) *WebKeyRegistry {
	client := &httputil.Client{Server: base}
	return &WebKeyRegistry{client: client}
}

// Keys returns the public keys of the given user.
func (s *WebKeyRegistry) Keys(user string) ([]*rsautil.PublicKey, error) {
	if !IsSimpleName(user) {
		return nil, errcode.InvalidArgf("unsupported user name: %q", user)
	}
	bs, err := s.client.GetBytes(user)
	if err != nil {
		return nil, err
	}
	return rsautil.ParsePublicKeys(bs)
}
