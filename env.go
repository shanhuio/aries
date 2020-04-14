package aries

import (
	"context"
	"net/http"
)

// Env provides the generic config structure for starting a service.
type Env struct {
	// Context is the main context for running the service.
	// This is often just context.Background()
	Context context.Context

	// Config to make the server.
	Config interface{}

	// For the server to send outgoing HTTP requests.
	Transport *http.Transport

	// If this is testing environment.
	Testing bool
}

// BuildFunc builds a service using the given config and logger.
type BuildFunc func(env *Env) (Service, error)
