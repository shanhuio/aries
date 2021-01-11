package oauth

import (
	"net/url"
	"strings"

	"shanhu.io/misc/errcode"
)

// DiscardURLServerParts discards the server parts of an URL,
// including scheme, host, port and user info.
func DiscardURLServerParts(u *url.URL) *url.URL {
	cp := *u
	cp.Scheme = ""
	cp.Opaque = ""
	cp.User = nil
	cp.Host = ""
	return &cp
}

// ParseRedirect parses an in-site redirection URL.
// The server parts (scheme, host, port, user info) are discarded.
func ParseRedirect(redirect string) (string, error) {
	if redirect == "" {
		return "", nil
	}

	u, err := url.Parse(redirect)
	if err != nil {
		return "", err
	}
	if !strings.HasPrefix(u.Path, "/") {
		return "", errcode.InvalidArgf(
			"redirect path part %q is not absolute", u.Path,
		)
	}

	return DiscardURLServerParts(u).String(), nil
}
