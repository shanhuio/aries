package aries

import (
	"net/http"
)

// Service is a interface similar to Func
type Service interface {
	Serve(c *C) error
}

// Serve wraps a service into an HTTP handler.
func Serve(s Service) http.Handler {
	f, ok := s.(Func)
	if ok {
		// if it is a function already, we do not need to do the wrapping.
		return f
	}
	return Func(s.Serve)
}
