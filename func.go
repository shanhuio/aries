package aries

import (
	"io"
	"net"
	"net/http"
)

// Func defines an HTTP handling function.
type Func func(c *C) error

// Serve implements the service interface.
func (f Func) Serve(c *C) error { return f(c) }

func (f Func) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := NewContext(w, req)
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

// ServeAt launches the handler as an HTTP service at the given
// listener.
func (f Func) ServeAt(lis net.Listener) error {
	s := &http.Server{Handler: f}
	return s.Serve(lis)
}

// StringFunc creates a Func that always reply the given string.
func StringFunc(s string) Func {
	return func(c *C) error {
		io.WriteString(c.Resp, s)
		return nil
	}
}
