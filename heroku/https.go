package heroku

import (
	"shanhu.io/aries"
)

// RedirectHTTPS redirects incoming HTTPS requests to HTTPS.
func RedirectHTTPS(c *aries.C) bool {
	if c.HTTPS {
		return false
	}

	u := c.Req.URL
	u.Host = c.Req.Host
	u.Scheme = "https"
	c.Redirect(u.String())
	return true
}
