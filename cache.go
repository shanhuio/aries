package aries

// NeverCache sets the Cache-Control header to "no-store".
func NeverCache(c *C) {
	c.Resp.Header().Set("Cache-Control", "max-age=0; no-store")
}
