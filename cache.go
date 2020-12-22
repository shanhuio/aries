package aries

// NoCache sets the Cache-Control header to "no-cache".
func NoCache(c *C) {
	c.Resp.Header().Set("Cache-Control", "no-cache")
}
