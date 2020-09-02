package oauth

import (
	"encoding/json"
	"fmt"

	"golang.org/x/oauth2"
	goauth2 "golang.org/x/oauth2/google"
	"shanhu.io/aries"
	"shanhu.io/misc/signer"
	"shanhu.io/misc/strutil"
)

// GoogleApp stores the configuration of a Google oauth2 application.
type GoogleApp struct {
	ID          string
	Secret      string
	RedirectURL string

	WithProfile bool
	Scopes      []string
}

const (
	googleEmailScope   = "https://www.googleapis.com/auth/userinfo.email"
	googleProfileScope = "https://www.googleapis.com/auth/userinfo.profile"
)

type google struct{ c *Client }

func newGoogle(app *GoogleApp, s *signer.Sessions) *google {
	scopeSet := make(map[string]bool)
	// Google OAuth has to have at least one scope to get user ID.
	scopeSet[googleEmailScope] = true
	if app.WithProfile {
		scopeSet[googleProfileScope] = true
	}
	scopes := strutil.SortedList(scopeSet)
	if scopes == nil {
		scopes = []string{}
	}
	c := NewClient(
		&oauth2.Config{
			ClientID:     app.ID,
			ClientSecret: app.Secret,
			Scopes:       scopes,
			Endpoint:     goauth2.Endpoint,
			RedirectURL:  app.RedirectURL,
		}, s, MethodGoogle,
	)
	return &google{c: c}
}

func (g *google) client() *Client { return g.c }

func (g *google) callback(c *aries.C) (*UserMeta, *State, error) {
	tok, state, err := g.c.TokenState(c)
	if err != nil {
		return nil, nil, err
	}

	const url = "https://www.googleapis.com/oauth2/v3/userinfo"
	bs, err := g.c.Get(c.Context, tok, url)
	if err != nil {
		return nil, nil, err
	}

	var user struct {
		Email string `json:"email"`
		Name  string `json:"name"`
	}
	if err := json.Unmarshal(bs, &user); err != nil {
		return nil, nil, err
	}
	email := user.Email
	if email == "" {
		return nil, nil, fmt.Errorf("empty login")
	}
	name := user.Name
	if name == "" {
		name = "no-name"
	}
	return &UserMeta{
		Method: MethodGoogle,
		ID:     email,
		Name:   name,
		Email:  email,
	}, state, nil
}
