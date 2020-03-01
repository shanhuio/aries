package oauth

import (
	"shanhu.io/aries"
)

type userMeta struct {
	id    string
	name  string
	email string
}

type metaExchange interface {
	callback(c *aries.C) (*userMeta, *State, error)
}
