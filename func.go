package aries

import (
	"fmt"
	"net/http"
)

// Func defines an HTTP handling function.
type Func func(c *C) error

func okHandler(c *C) error {
	fmt.Fprint(c.Resp, "ok")
	return nil
}

func (f Func) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := NewContext(w, req)
	if c.Path == "/ok" {
		c.ErrCode(okHandler(c))
		return
	}
	c.ErrCode(f(c))
}

// ListenAndServe launches the handler as an HTTP service.
func (f Func) ListenAndServe(addr string) error {
	s := &http.Server{
		Addr:    addr,
		Handler: f,
	}
	return s.ListenAndServe()
}

// ListenAndServe launches the handler as an HTTP service.
func ListenAndServe(addr string, f Func) error {
	return f.ListenAndServe(addr)
}
