package oauth

import (
	"time"

	"shanhu.io/aries"
	"shanhu.io/misc/signer"
)

// Module is a module that handles stuff related to oauth.
type Module struct {
	c        *Config
	github   *github
	google   *google
	sessions *signer.Sessions
	redirect string

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
	return ret
}

// Mux makes a webapp module mux that executes the oauth flow on the server
// side.
func (mod *Module) Mux() *aries.Mux {
	m := aries.NewMux()
	if mod.c.ByPass != "" {
		m.Exact("/signin-bypass", func(c *aries.C) {
			mod.SetupCookie(c, mod.c.ByPass)
			mod.signInRedirect(c)
		})
	}

	m.Exact("/signout", func(c *aries.C) {
		c.ClearCookie("session")
		c.Redirect(mod.redirect)
	})

	if mod.c.KeyStore != nil {
		m.Exact("/pubkey/signin", func(c *aries.C) {
			req := new(LoginRequest)
			if err := aries.UnmarshalJSONBody(c, req); c.Error(500, err) {
				return
			}

			keyBytes, err := mod.c.KeyStore.Key(req.User)
			if c.Error(500, err) {
				return
			}

			k, err := unmarshalPublicKey(keyBytes)
			if c.Error(500, err) {
				return
			}

			const loginCredsWindow = time.Minute * 5
			sig := signer.NewRSATimeSigner(k, loginCredsWindow)
			if c.Error(403, sig.Check(req.SignedTime)) {
				return
			}

			session, expires := mod.NewCreds(req.User)
			resp := &Creds{
				User:    req.User,
				Token:   session,
				Expires: expires.UnixNano(),
			}
			aries.ReplyJSON(c, resp)
		})
	}

	if mod.github != nil {
		m.Exact("/github/signin", func(c *aries.C) {
			c.Redirect(mod.github.signInURL())
		})

		m.Exact("/github/callback", func(c *aries.C) {
			user, err := mod.github.callback(c.Req)
			if err != nil {
				c.AltError(err, 500, "callback failed")
				return
			}
			mod.signIn(c, "github", user)
		})
	}

	if mod.google != nil {
		m.Exact("/google/signin", func(c *aries.C) {
			c.Redirect(mod.google.signInURL())
		})

		m.Exact("/google/callback", func(c *aries.C) {
			user, err := mod.google.callback(c.Req)
			if err != nil {
				c.AltError(err, 500, "callback failed")
				return
			}
			mod.signIn(c, "google", user)
		})
	}

	return m
}

func readSessionToken(c *aries.C) (string, string) {
	if bearer := aries.Bearer(c); bearer != "" {
		return bearer, "bearer"
	}
	return c.ReadCookie("session"), "cookie"
}

// SetupCookie sets up the cookie for a particular user.
func (mod *Module) SetupCookie(c *aries.C, user string) {
	session, expires := mod.sessions.New([]byte(user))
	c.WriteCookie("session", session, expires)
}

func (mod *Module) signInRedirect(c *aries.C) {
	c.Redirect(mod.redirect)
}

func (mod *Module) signIn(c *aries.C, method, user string) {
	if mod.c.LoginCheck == nil {
		return
	}

	id := mod.c.LoginCheck(c, method, user)
	if id != "" {
		mod.SetupCookie(c, id)
		mod.signInRedirect(c)
	}
}

func (mod *Module) checkUser(c *aries.C) (valid, needRefresh bool) {
	c.User = ""

	session, method := readSessionToken(c)
	bs, left, ok := mod.sessions.Check(session)
	if !ok {
		return false, false
	}
	needRefresh = left < mod.sessionRefresh && method == "cookie"

	user := string(bs)
	if mod.c.Check == nil {
		c.User = user
		return true, needRefresh
	}

	u, ok := mod.c.Check(user)
	if !ok {
		return false, false
	}

	c.User = user
	c.Data["user"] = u
	return true, needRefresh
}

// NewCreds creates new credentials for the user.
func (mod *Module) NewCreds(user string) (string, time.Time) {
	return mod.sessions.New([]byte(user))
}

// Check checks the user credentials.
func (mod *Module) Check(c *aries.C) bool {
	ok, needRefresh := mod.checkUser(c)
	if !ok {
		c.ClearCookie("session")
	} else if needRefresh {
		mod.SetupCookie(c, c.User)
	}
	return ok
}
