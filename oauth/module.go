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

type service struct {
	r *aries.Router
	m *Module
}

func (s *service) Setup(c *aries.C) error {
	_, err := s.m.Check(c)
	return err
}

func (s *service) Serve(c *aries.C) error { return s.r.Serve(c) }

// Auth makes a aries.Auth that executes the oauth flow on the server side.
func (mod *Module) Auth() aries.Auth {
	r := aries.NewRouter()
	if mod.c.ByPass != "" {
		r.File("signin-bypass", func(c *aries.C) error {
			mod.SetupCookie(c, mod.c.ByPass)
			c.Redirect(mod.redirect)
			return nil
		})
	}

	r.File("signout", func(c *aries.C) error {
		c.ClearCookie("session")
		c.Redirect(mod.redirect)
		return nil
	})

	if mod.c.KeyStore != nil {
		r.File("pubkey/signin", func(c *aries.C) error {
			req := new(LoginRequest)
			if err := aries.UnmarshalJSONBody(c, req); err != nil {
				return err
			}

			keyBytes, err := mod.c.KeyStore.Key(req.User)
			if err != nil {
				return err
			}

			k, err := rsautil.ParsePublicKey(keyBytes)
			if err != nil {
				return err
			}

			const loginCredsWindow = time.Minute * 5
			sig := signer.NewRSATimeSigner(k, loginCredsWindow)
			if err := sig.Check(req.SignedTime); err != nil {
				return errcode.Add(errcode.Unauthorized, err)
			}

			session, expires := mod.NewCreds(req.User)
			resp := &Creds{
				User:    req.User,
				Token:   session,
				Expires: expires.UnixNano(),
			}
			return aries.ReplyJSON(c, resp)
		})
	}

	if mod.github != nil {
		r.File("github/signin", func(c *aries.C) error {
			state := &State{Dest: mod.redirect}
			c.Redirect(mod.github.client.SignInURL(state))
			return nil
		})

		r.File("github/callback", func(c *aries.C) error {
			user, state, err := mod.github.callback(c)
			if err != nil {
				log.Println("github callback: ", err)
				return errcode.Internalf("callback failed")
			}
			return mod.signIn(c, "github", user, state.Dest)
		})
	}

	if mod.google != nil {
		r.File("google/signin", func(c *aries.C) error {
			state := &State{Dest: mod.redirect}
			c.Redirect(mod.google.client.SignInURL(state))
			return nil
		})

		r.File("google/callback", func(c *aries.C) error {
			user, state, err := mod.google.callback(c)
			if err != nil {
				log.Println("google callback: ", err)
				return errcode.Internalf("callback failed")
			}
			return mod.signIn(c, "google", user, state.Dest)
		})
	}

	return &service{
		m: mod,
		r: r,
	}
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

func (mod *Module) signIn(c *aries.C, method, user, dest string) error {
	if mod.c.LoginCheck == nil {
		return nil
	}

	id, err := mod.c.LoginCheck(c, method, user)
	if err != nil {
		return err
	}
	if id != "" {
		mod.SetupCookie(c, id)
		c.Redirect(dest)
	}
	return nil
}

func (mod *Module) checkUser(c *aries.C) (
	valid, needRefresh bool, err error,
) {
	c.User = ""

	session, method := readSessionToken(c)
	bs, left, ok := mod.sessions.Check(session)
	if !ok {
		return false, false, nil
	}
	needRefresh = left < mod.sessionRefresh && method == "cookie"

	user := string(bs)
	if mod.c.Check == nil {
		c.User = user
		c.UserLevel = 0
		return true, needRefresh, nil
	}

	u, lvl, err := mod.c.Check(user)
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

// NewCreds creates new credentials for the user.
func (mod *Module) NewCreds(user string) (string, time.Time) {
	return mod.sessions.New([]byte(user))
}

// Check checks the user credentials.
func (mod *Module) Check(c *aries.C) (bool, error) {
	ok, needRefresh, err := mod.checkUser(c)
	if err != nil {
		return false, err
	}
	if !ok {
		c.ClearCookie("session")
	} else if needRefresh {
		mod.SetupCookie(c, c.User)
	}
	return ok, nil
}

// AuthSetup setups the user authorization context.
func (mod *Module) AuthSetup(c *aries.C) { mod.Check(c) }
