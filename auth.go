package aries

// Auth is an authentication service.
type Auth interface {
	Service

	// Setup sets up the authentication in context.
	Setup(c *C) error
}
