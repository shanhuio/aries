package aries

// NeverCache sets the Cache-Control header to "max-age=0; no-store".
func NeverCache(c *C) {
	c.Resp.Header().Set("Cache-Control", "max-age=0; no-store")
}
