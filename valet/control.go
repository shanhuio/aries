package valet

import (
	"fmt"
	"path/filepath"
	"time"

	"shanhu.io/aries"
	"shanhu.io/aries/oauth"
	"shanhu.io/aries/sitter"
	"shanhu.io/misc/errcode"
)

// the controller server
type control struct {
	c        *Config
	oauth    *oauth.Module
	oauthMux *aries.Mux
	apiMux   *aries.Mux
}

func (s *control) f(f func(s *control, c *aries.C) error) aries.Func {
	return func(c *aries.C) error { return f(s, c) }
}

func serveIndex(s *control, c *aries.C) error {
	fmt.Fprint(c.Resp, "controller")
	return nil
}

func serveDeploy(s *control, c *aries.C) error {
	var req struct {
		Name string
	}

	err := aries.UnmarshalJSONBody(c, &req)
	if err != nil {
		return err
	}

	dir := filepath.Join("/prod", req.Name, "pkg")
	server := fmt.Sprintf("%s:8105", req.Name)
	if err := sitter.Push(dir, server, c.Resp); err != nil {
		fmt.Fprintf(c.Resp, "error: %s\n", err)
	}
	return nil
}

func makeAPIMux(s *control) *aries.Mux {
	m := aries.NewMux()
	m.Exact("/", s.f(serveIndex))
	m.Exact("/deploy", s.f(serveDeploy))
	return m
}

func makeControl(c *Config) aries.Func {
	s := &control{
		c: c,
	}

	s.oauth = oauth.NewModule(&oauth.Config{
		KeyStore: oauth.NewFileKeyStore(
			map[string]string{
				c.Admin: c.PublicKey,
			},
		),
		SessionLifeTime: time.Hour * 24 * 3,
		SessionKey:      []byte(c.SessionKey),
	})
	s.oauthMux = s.oauth.Mux()
	s.apiMux = makeAPIMux(s)

	return s.f(serveControl)
}

func serveControl(s *control, c *aries.C) error {
	if hit, err := s.oauthMux.Serve(c); hit {
		return err
	}
	s.oauth.Check(c)

	if c.User == "" {
		fmt.Fprintln(c.Resp, randomMotto())
		return nil
	}

	if hit, err := s.apiMux.Serve(c); hit {
		return err
	}

	return errcode.NotFoundf("nothing here")
}
