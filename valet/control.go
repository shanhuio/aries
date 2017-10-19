package valet

import (
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"shanhu.io/aries"
	"shanhu.io/aries/oauth"
	"shanhu.io/aries/sitter"
)

// the controller server
type control struct {
	c        *Config
	oauth    *oauth.Module
	oauthMux *aries.Mux
	apiMux   *aries.Mux
}

func (s *control) f(f func(s *control, c *aries.C)) func(c *aries.C) {
	return func(c *aries.C) { f(s, c) }
}

func serveIndex(s *control, c *aries.C) {
	fmt.Fprint(c.Resp, "controller")
}

func serveDeploy(s *control, c *aries.C) {
	var req struct {
		Name string
	}

	err := aries.UnmarshalJSONBody(c, &req)
	if c.Error(400, err) {
		return
	}

	dir := filepath.Join("/prod", req.Name, "pkg")
	server := fmt.Sprintf("%s:8105", req.Name)
	if err := sitter.Push(dir, server, c.Resp); err != nil {
		fmt.Fprintf(c.Resp, "error: %s\n", err)
	}
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

func serveControl(s *control, c *aries.C) {
	if s.oauthMux.Serve(c) {
		return
	}
	s.oauth.Check(c)

	if c.User == "" {
		fmt.Fprintln(c.Resp, randomMotto())
		return
	}

	if s.apiMux.Serve(c) {
		return
	}

	http.Error(c.Resp, "nothing found here", 404)
}
