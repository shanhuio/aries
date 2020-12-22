package oauth

import (
	"net/url"
	"strings"

	"shanhu.io/misc/errcode"
)

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

	cp := &url.URL{
		Path:        u.Path,
		RawPath:     u.RawPath,
		ForceQuery:  u.ForceQuery,
		RawQuery:    u.RawQuery,
		Fragment:    u.Fragment,
		RawFragment: u.RawFragment,
	}
	return cp.String(), nil
}
