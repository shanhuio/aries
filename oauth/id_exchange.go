package oauth

import (
	"shanhu.io/aries"
)

type idExchanger interface {
	callback(c *aries.C) (string, *State, error)
}
