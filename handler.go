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
	h.Func(&C{
		Path:  req.URL.Path,
		Resp:  w,
		Req:   req,
		HTTPS: h.HTTPS,
		Data:  make(map[string]interface{}),
	})
}

// ListenAndServe launches the handler as an HTTP service.
func (h *Handler) ListenAndServe(addr string) error {
	s := &http.Server{
		Addr:    addr,
		Handler: h,
	}
	return s.ListenAndServe()
}

// HandleFunc a context serving function into a context serving function.
func HandleFunc(f Func, https bool) http.HandlerFunc {
	h := &Handler{
		Func:  f,
		HTTPS: https,
	}
	return h.ServeHTTP
}
