package oauth

import (
	"shanhu.io/aries"
)

// UserMeta returns the meta data returned by a sign in.
type UserMeta struct {
	Method string
	ID     string
	Name   string // Screen name.
	Email  string
}

type metaExchange interface {
	callback(c *aries.C) (*UserMeta, *State, error)
}

type provider interface {
	metaExchange
	client() *Client
}
