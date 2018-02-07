package aries

import (
	"shanhu.io/misc/errcode"
)

// MuxSet is a set of muxes that
type MuxSet struct {
	Auth      *Mux
	AuthSetup func(c *C)

	Resource *Mux
	Guest    *Mux
	Admin    *Mux

	InternalSignIn Func
}

func serveMux(m *Mux, c *C) error {
	if m == nil {
		return Miss
	}
	return m.Serve(c)
}

func isAdmin(c *C) bool {
	return c.User != "" && c.UserLevel > 0
}

// Serve serves the incoming request with the mux set.
func (s *MuxSet) Serve(c *C) error {
	if err := serveMux(s.Auth, c); err != Miss {
		return err
	}
	if s.AuthSetup != nil {
		s.AuthSetup(c)
	}

	if err := serveMux(s.Resource, c); err != Miss {
		return err
	}
	if err := serveMux(s.Guest, c); err != Miss {
		return err
	}

	if isAdmin(c) {
		if err := serveMux(s.Admin, c); err != Miss {
			return err
		}
	}

	return Miss
}

// ServeInternal serves the incoming request with the mux set, but only serves
// resource for normal users, and allows only admins (users with positive
// level) to visit the guest mux.
func (s *MuxSet) ServeInternal(c *C) error {
	if err := serveMux(s.Auth, c); err != Miss {
		return err
	}
	if s.AuthSetup != nil {
		s.AuthSetup(c)
	}

	if err := serveMux(s.Resource, c); err != Miss {
		return err
	}

	if !isAdmin(c) {
		if c.Path == "/" {
			if s.InternalSignIn != nil {
				return s.InternalSignIn(c)
			}
			return errcode.Unauthorizedf("please sign in")
		}
		c.Redirect("/")
		return nil
	}

	if err := serveMux(s.Guest, c); err != Miss {
		return err
	}
	if err := serveMux(s.Admin, c); err != Miss {
		return err
	}

	return Miss
}
