package aries

// ServiceSet is a set of muxes that
type ServiceSet struct {
	Auth Auth

	Resource   Service
	Guest      Service
	User       Service
	Admin      Service
	IsInternal func(c *C) bool

	InternalSignIn Func
}

func serveService(m Service, c *C) error {
	if m == nil {
		return Miss
	}
	return m.Serve(c)
}

func (s *ServiceSet) isInternal(c *C) bool {
	if s.IsInternal == nil {
		return c.User != "" && c.UserLevel > 0
	}
	return s.IsInternal(c)
}

func (s *ServiceSet) serveAuth(c *C) error {
	if s.Auth == nil {
		return nil
	}
	if err := s.Auth.Serve(c); err != Miss {
		return err
	}
	return s.Auth.Setup(c)
}

// Serve serves the incoming request with the mux set.
func (s *ServiceSet) Serve(c *C) error {
	if err := s.serveAuth(c); err != nil {
		return err
	}

	if err := serveService(s.Resource, c); err != Miss {
		return err
	}
	if err := serveService(s.Guest, c); err != Miss {
		return err
	}
	if c.User != "" {
		if err := serveService(s.User, c); err != Miss {
			return err
		}
	}
	if s.isInternal(c) {
		if err := serveService(s.Admin, c); err != Miss {
			return err
		}
	}

	return Miss
}

// ServeInternal serves the incoming request with the mux set, but only serves
// resource for normal users, and allows only admins (users with positive
// level) to visit the guest mux.
func (s *ServiceSet) ServeInternal(c *C) error {
	if err := serveService(s.Auth, c); err != Miss {
		return err
	}
	if s.Auth != nil {
		if err := s.Auth.Setup(c); err != nil {
			return err
		}
	}

	if err := serveService(s.Resource, c); err != Miss {
		return err
	}

	if !s.isInternal(c) {
		if c.Path == "/" {
			if s.InternalSignIn != nil {
				return s.InternalSignIn(c)
			}
			return NeedSignIn
		}
		c.Redirect("/")
		return nil
	}

	if err := serveService(s.Guest, c); err != Miss {
		return err
	}
	if err := serveService(s.User, c); err != Miss {
		return err
	}
	if err := serveService(s.Admin, c); err != Miss {
		return err
	}

	return Miss
}
