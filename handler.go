package aries

import (
	"net/http"
)

// Handler implements the standard http interface.
type Handler struct {
	Func
	HTTPS bool
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := NewContext(w, req, h.HTTPS)
	c.ErrCode(h.Func(c))
}

// ListenAndServe launches the handler as an HTTP service.
func (h *Handler) ListenAndServe(addr string) error {
	s := &http.Server{
		Addr:    addr,
		Handler: h,
	}
	return s.ListenAndServe()
}

// HandlerFunc wraps a context serving function into an HTTP handler function.
func HandlerFunc(f Func, https bool) http.HandlerFunc {
	h := &Handler{
		Func:  f,
		HTTPS: https,
	}
	return h.ServeHTTP
}
