package aries

import (
	"net/http"
)

// Env provides the generic config structure for starting a service.
type Env struct {
	// Config to make the server.
	Config interface{}

	// For the server to log stuff.
	Logger *Logger

	// For the server to send outgoing HTTP requests.
	Transport http.RoundTripper
}

// BuildFunc builds a service using the given config and logger.
type BuildFunc func(env *Env) (Service, error)
