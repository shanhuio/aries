package oauth

import (
	"shanhu.io/aries"
)

// SignUp is an HTTP module that handles user signups.
type SignUp struct {
	redirect string
	purpose  string
	module   *Module
	router   *aries.Router
}

// SignUpConfig is the config for creating a signup module.
type SignUpConfig struct {
	Redirect string
}

// NewSignUp creates a new sign up module.
func NewSignUp(m *Module, c *SignUpConfig) *SignUp {
	s := &SignUp{
		purpose: "signup",
		redirect: c.Redirect,
		module:   m,
	}

	s.router = s.makeRouter()
	return s
}

// Serve serves the incoming HTTP request.
func (s *SignUp) Serve(c *aries.C) error {
	return s.router.Serve(c)
}

func (s *SignUp) makeRouter() *aries.Router {
	r := aries.NewRouter()
	methods := s.module.Methods()
	for _, m := range methods {
		r.File(m, s.handler(m))
	}
	return r
}

// Purpose returns the purpose string.
func (s *SignUp) Purpose() string {
	return s.purpose
}

func (s *SignUp) handler(m string) aries.Func {
	return func(c *aries.C) error {
		state := &State{
			Dest:     s.redirect,
			NoCookie: true,
			Purpose:  s.purpose,
		}
		s.module.SignIn(c, m, state)
		return nil
	}
}
