package aries

import (
	"fmt"
	"net/http"
)

// StaticFiles is a module that serves static files.
type StaticFiles struct {
	cacheControl string
	h            http.Handler
}

func cacheControl(ageSecs int) string {
	return fmt.Sprintf("max-age=%d; must-revalidate", ageSecs)
}

// NewStaticFiles creates a module that serves static files.
func NewStaticFiles(p string) *StaticFiles {
	return &StaticFiles{
		cacheControl: cacheControl(10),
		h:            http.FileServer(http.Dir(p)),
	}
}

// CacheAge sets the maximum age for the cache.
func (s *StaticFiles) CacheAge(ageSecs int) {
	if ageSecs < 0 {
		s.cacheControl = ""
	} else {
		s.cacheControl = cacheControl(ageSecs)
	}
}

// Serve serves incoming HTTP requests.
func (s *StaticFiles) Serve(c *C) error {
	c.Req.URL.Path = c.Path
	if s.cacheControl != "" {
		c.Resp.Header().Add("Cache-Control", s.cacheControl)
	}
	s.h.ServeHTTP(c.Resp, c.Req)
	return nil
}
