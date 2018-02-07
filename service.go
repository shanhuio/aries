package aries

// Service is a interface similar to Func
type Service interface {
	Serve(c *C) error
}
