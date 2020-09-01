package oauth

import (
	"log"
	"time"

	"shanhu.io/aries"
	"shanhu.io/misc/errcode"
	"shanhu.io/misc/rsautil"
	"shanhu.io/misc/signer"
)

// Module is a module that handles stuff related to oauth.
type Module struct {
	c            *Config
	github       *github
	google       *google
	digitalOcean *digitalOcean
	sessions     *signer.Sessions
	redirect     string

	sessionRefresh time.Duration
}

// NewModule creates a new oauth module with the given config.
func NewModule(c *Config) *Module {
	redirect := c.Redirect
	if redirect == "" {
		redirect = "/"
	}

	sessionLifeTime := c.SessionLifeTime
	if sessionLifeTime <= 0 {
		sessionLifeTime = time.Hour * 24 * 7 // roughly a week
	}
	sessionRefresh := c.SessionRefresh
	if sessionRefresh <= 0 || sessionRefresh > sessionLifeTime {
		sessionRefresh = sessionLifeTime / 5 * 4
	}

	ret := &Module{
		c: c,
		sessions: signer.NewSessions(
			c.SessionKey,
			sessionLifeTime,
		),
		redirect:       redirect,
		sessionRefresh: sessionRefresh,
	}

	const ttl time.Duration = time.Hour
	states := signer.NewSessions(c.StateKey, ttl)

	if c.GitHub != nil {
		ret.github = newGitHub(c.GitHub, states)
	}
	if c.Google != nil {
		ret.google = newGoogle(c.Google, states)
	}
	if c.DigitalOcean != nil {
		ret.digitalOcean = newDigitalOcean(c.DigitalOcean, states)
	}
	return ret
}

type service struct {
	r *aries.Router
	m *Module
}

func (s *service) Setup(c *aries.C) error {
	_, err := s.m.Check(c)
	return err
}

func (s *service) Serve(c *aries.C) error { return s.r.Serve(c) }

func (m *Module) pubKeySignIn(c *aries.C, r *LoginRequest) (*Creds, error) {
	if r.SignedTime == nil {
		return nil, errcode.InvalidArgf("signature missing")
	}

	keys, err := m.c.KeyStore.Keys(r.User)
	if err != nil {
		return nil, err
	}

	var key *rsautil.PublicKey
	for _, k := range keys {
		if k.HashStr() == r.SignedTime.KeyID {
			key = k
			break
		}
	}
	if key == nil {
		return nil, errcode.Unauthorizedf("signing key not authorized")
	}

	const window = time.Minute * 5
	if err := signer.CheckRSATimeSignature(
		r.SignedTime, key.Key(), window,
	); err != nil {
		return nil, errcode.Add(errcode.Unauthorized, err)
	}

	token, expires := m.newCreds(r.User, time.Duration(r.TTL))
	return &Creds{
		User:    r.User,
		Token:   token,
		Expires: expires.UnixNano(),
	}, nil
}

// Auth makes a aries.Auth that executes the oauth flow on the server side.
func (m *Module) Auth() aries.Auth {
	r := aries.NewRouter()
	if m.c.Bypass != "" {
		r.File("signin-bypass", func(c *aries.C) error {
			m.SetupCookie(c, m.c.Bypass)
			c.Redirect(m.redirect)
			return nil
		})
	}
	r.File("signout", func(c *aries.C) error {
		c.ClearCookie("session")
		c.Redirect(m.redirect)
		return nil
	})

	if m.c.KeyStore != nil {
		r.JSONCallMust("pubkey/signin", m.pubKeySignIn)
	}

	if g := m.github; g != nil {
		r.File("github/signin", m.signInHandler(g.client()))
		r.File("github/callback", m.callbackHandler("github", g))
	}
	if g := m.google; g != nil {
		r.File("google/signin", m.signInHandler(g.client()))
		r.File("google/callback", m.callbackHandler("google", g))
	}
	if do := m.digitalOcean; do != nil {
		r.File("digitalocean/signin", m.signInHandler(do.client()))
		r.File("digitalocean/callback", m.callbackHandler(
			"digitalocean", do,
		))
	}

	return &service{m: m, r: r}
}

func readSessionToken(c *aries.C) (string, string) {
	if bearer := aries.Bearer(c); bearer != "" {
		return bearer, "bearer"
	}
	return c.ReadCookie("session"), "cookie"
}

// SetupCookie sets up the cookie for a particular user.
func (m *Module) SetupCookie(c *aries.C, user string) {
	token, expires := m.newCreds(user, 0)
	c.WriteCookie("session", token, expires)
}

func (m *Module) signIn(
	c *aries.C, user *UserMeta, state *State,
) error {
	if m.c.LoginCheck == nil {
		if !state.NoCookie {
			m.SetupCookie(c, user.ID)
		}
		c.Redirect(state.Dest)
		return nil
	}

	id, err := m.c.LoginCheck(c, user)
	if err != nil {
		return err
	}
	if id != "" {
		if !state.NoCookie {
			m.SetupCookie(c, id)
		}
		c.Redirect(state.Dest)
	}
	return nil
}

func (m *Module) checkUser(c *aries.C) (
	valid, needRefresh bool, err error,
) {
	c.User = ""

	session, method := readSessionToken(c)
	bs, left, ok := m.sessions.Check(session)
	if !ok {
		return false, false, nil
	}
	needRefresh = left < m.sessionRefresh && method == "cookie"

	user := string(bs)
	if m.c.Check == nil {
		c.User = user
		c.UserLevel = 0
		if user != "" {
			c.UserLevel = 1
		}
		return true, needRefresh, nil
	}

	u, lvl, err := m.c.Check(user)
	if err != nil {
		return false, false, err
	}
	if lvl < 0 {
		return false, false, nil
	}

	c.User = user
	c.UserLevel = lvl
	if u != nil {
		c.Data["user"] = u
	}
	return true, needRefresh, nil
}

func (m *Module) newCreds(user string, ttl time.Duration) (string, time.Time) {
	return m.sessions.New([]byte(user), ttl)
}

// Check checks the user credentials.
func (m *Module) Check(c *aries.C) (bool, error) {
	ok, needRefresh, err := m.checkUser(c)
	if err != nil {
		return false, err
	}
	if !ok {
		c.ClearCookie("session")
	} else if needRefresh {
		m.SetupCookie(c, c.User)
	}
	return ok, nil
}

// AuthSetup setups the user authorization context.
func (m *Module) AuthSetup(c *aries.C) { m.Check(c) }

func (m *Module) signInHandler(client *Client) aries.Func {
	return func(c *aries.C) error {
		redirect := c.Req.URL.Query().Get("redirect")
		if redirect == "" {
			redirect = m.redirect
		}
		state := &State{Dest: redirect}
		if c.Req.URL.Query().Get("cookie") == "false" {
			state.NoCookie = true
		}
		c.Redirect(client.SignInURL(state))
		return nil
	}
}

func (m *Module) callbackHandler(method string, x metaExchange) aries.Func {
	return func(c *aries.C) error {
		user, state, err := x.callback(c)
		if err != nil {
			log.Printf("%s callback: %s", method, err)
			return errcode.Internalf("%s callback failed", method)
		}
		if user == nil {
			return errcode.Internalf(
				"%s callback: get user info failed", method,
			)
		}

		return m.signIn(c, user, state)
	}
}
