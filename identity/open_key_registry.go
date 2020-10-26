package identity

import (
	"net/url"

	"shanhu.io/misc/errcode"
)

// OpenKeyRegistry connects to a keystore based on the given URL string.
func OpenKeyRegistry(urlStr string) (KeyRegistry, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	switch u.Scheme {
	case "http", "https":
		u, err := url.Parse(urlStr)
		if err != nil {
			return nil, err
		}
		return NewWebKeyRegistry(u), nil
	case "file", "":
		r, err := NewDirKeyRegistry(u.Path)
		if err != nil {
			return nil, err
		}
		return r, nil
	}
	return nil, errcode.InvalidArgf("unsupported url scheme: %q", u.Scheme)
}
