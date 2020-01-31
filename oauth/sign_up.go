package oauth

import (
	"log"
	"time"

	"shanhu.io/aries"
	"shanhu.io/misc/errcode"
	"shanhu.io/misc/signer"
)

// SignUpRequest contains information of a sign up request.
type SignUpRequest struct {
	Method string
	ID     string
	Email  string
}

// SignUp is an HTTP module that handles user signups.
type SignUp struct {
	google     *google
	router     *aries.Router
	reqHandler func(c *aries.C, req *SignUpRequest) error
}

// SignUpConfig is the config for creating a signup module.
type SignUpConfig struct {
	StateKey []byte
	Google   *GoogleApp

	HandleRequest func(c *aries.C, req *SignUpRequest) error
}

// NewSignUp creates a new sign up module.
func NewSignUp(c *SignUpConfig) *SignUp {
	s := &SignUp{
		reqHandler: c.HandleRequest,
	}

	const ttl time.Duration = time.Hour
	states := signer.NewSessions(c.StateKey, ttl)

	if c.Google != nil {
		s.google = newGoogle(c.Google, states)
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

	if g := s.google; g != nil {
		r.File("google", s.handler(g.client()))
		r.File("google:callback", s.callback("google", g))
	}

	return r
}

func (s *SignUp) handler(client *Client) aries.Func {
	return func(c *aries.C) error {
		state := new(State)
		c.Redirect(client.SignInURL(state))
		return nil
	}
}

func (s *SignUp) callback(method string, x metaExchange) aries.Func {
	return func(c *aries.C) error {
		user, _, err := x.callback(c)
		if err != nil {
			log.Printf("%s callback: %s", method, err)
			return errcode.Internalf("%s callback failed", method)
		}
		if user == nil {
			return errcode.Internalf(
				"%s callback: get user email failed", method,
			)
		}

		req := &SignUpRequest{
			Method: method,
			ID:     user.id,
			Email:  user.email,
		}
		return s.reqHandler(c, req)
	}
}
