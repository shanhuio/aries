package oauth

import (
	"encoding/json"
	"fmt"

	"golang.org/x/oauth2"
	goauth2 "golang.org/x/oauth2/google"
	"shanhu.io/aries"
	"shanhu.io/misc/signer"
)

// GoogleApp stores the configuration of a Google App.
type GoogleApp struct {
	ID          string
	Secret      string
	RedirectURL string
}

type google struct{ client *Client }

func newGoogle(app *GoogleApp, s *signer.Sessions) *google {
	c := NewClient(
		&oauth2.Config{
			ClientID:     app.ID,
			ClientSecret: app.Secret,
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
			},
			Endpoint:    goauth2.Endpoint,
			RedirectURL: app.RedirectURL,
		}, s,
	)
	return &google{c}
}

func (g *google) callback(c *aries.C) (string, *State, error) {
	tok, state, err := g.client.TokenState(c)
	if err != nil {
		return "", nil, err
	}

	const url = "https://www.googleapis.com/oauth2/v3/userinfo"
	bs, err := g.client.Get(c.Context, tok, url)
	if err != nil {
		return "", nil, err
	}

	var user struct {
		Email string `json:"email"`
	}
	if err := json.Unmarshal(bs, &user); err != nil {
		return "", nil, err
	}
	ret := user.Email
	if ret == "" {
		return "", nil, fmt.Errorf("empty login")
	}
	return ret, state, nil
}
