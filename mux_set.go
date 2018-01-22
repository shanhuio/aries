package aries

import (
	"shanhu.io/misc/errcode"
)

// MuxSet is a set of muxes that
type MuxSet struct {
	Auth      *Mux
	AuthCheck func(c *C)

	Resource *Mux
	Guest    *Mux
	Admin    *Mux

	InternalSignIn Func
}

func serveMux(m *Mux, c *C) (bool, error) {
	if m == nil {
		return false, nil
	}
	return m.Serve(c)
}

func isAdmin(c *C) bool {
	return c.User != "" && c.UserLevel > 0
}

// Serve serves the incoming request with the mux set.
func (s *MuxSet) Serve(c *C) (bool, error) {
	if hit, err := serveMux(s.Auth, c); hit {
		return true, err
	}
	if s.AuthCheck != nil {
		s.AuthCheck(c)
	}

	if hit, err := serveMux(s.Resource, c); hit {
		return true, err
	}
	if hit, err := serveMux(s.Guest, c); hit {
		return true, err
	}

	if isAdmin(c) {
		if hit, err := serveMux(s.Admin, c); hit {
			return true, err
		}
	}

	return false, nil
}

// ServeInternal serves the incoming request with the mux set, but only serves
// resource for normal users, and allows only admins (users with positive
// level) to visit the guest mux.
func (s *MuxSet) ServeInternal(c *C, signIn Func) (bool, error) {
	if hit, err := serveMux(s.Auth, c); hit {
		return true, err
	}
	if s.AuthCheck != nil {
		s.AuthCheck(c)
	}

	if hit, err := serveMux(s.Resource, c); hit {
		return true, err
	}

	if !isAdmin(c) {
		if c.Path == "/" {
			if s.InternalSignIn != nil {
				return true, s.InternalSignIn(c)
			}
			return true, errcode.Unauthorizedf("please sign in")
		}
		c.Redirect("/")
		return true, nil
	}

	if hit, err := serveMux(s.Guest, c); hit {
		return true, err
	}
	if hit, err := serveMux(s.Admin, c); hit {
		return true, err
	}

	return false, nil
}
