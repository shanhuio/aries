package heroku

import (
	"shanhu.io/aries"
)

// IsHTTPS checks if an incoming Heroku HTTP request is using HTTPS.
func IsHTTPS(c *aries.C) bool {
	return c.Req.Header.Get("X-Forwarded-Proto") == "https"
}

// RedirectHTTPS redirects incoming HTTPS requests to HTTPS.
func RedirectHTTPS(c *aries.C) bool {
	if IsHTTPS(c) {
		return false
	}

	u := c.Req.URL
	u.Host = c.Req.Host
	u.Scheme = "https"
	c.Redirect(u.String())
	return true
}
