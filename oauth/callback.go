package oauth

import (
	"shanhu.io/aries"
)

type idExchange interface {
	callback(c *aries.C) (string, *State, error)
}
