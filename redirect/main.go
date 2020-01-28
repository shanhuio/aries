package redirect

import (
	"shanhu.io/aries"
)

type config struct {
	RedirectToDomain string
}

type server struct {
	c *config
}

func (s *server) redirect(c *aries.C) error {
	u := *c.Req.URL // make a shallow copy
	u.Scheme = "https"
	u.Host = s.c.RedirectToDomain
	c.Redirect(u.String())
	return nil
}

func newServer(c *config) (aries.Func, error) {
	s := &server{c: c}
	return s.redirect, nil
}

func makeService(env *aries.Env) (aries.Service, error) {
	s, err := newServer(env.Config.(*config))
	if err != nil {
		return nil, err
	}
	return s, nil
}

// Main is the main entrance for the redirect service.
func Main() {
	aries.Main(makeService, new(config), "localhost:8000")
}
