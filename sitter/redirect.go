package sitter

import (
	"net/http"
	"net/http/httputil"
	"sync"
	"time"
)

type redirect struct {
	mu   sync.RWMutex
	host string
	*httputil.ReverseProxy
}

func newRedirect(host string) *redirect {
	s := &redirect{host: host}

	proxy := new(httputil.ReverseProxy)
	proxy.FlushInterval = time.Second * 3
	proxy.Director = func(req *http.Request) {
		s.mu.RLock()
		defer s.mu.RUnlock()

		url := req.URL
		url.Scheme = "http"
		url.Host = s.host
	}

	s.ReverseProxy = proxy
	return s
}

func (r *redirect) setHost(host string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.host = host
}
