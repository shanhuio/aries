package oauth

import (
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
	goauth2 "golang.org/x/oauth2/google"

	"shanhu.io/misc/signer"
)

// GoogleApp stores the configuration of a Google App.
type GoogleApp struct {
	ID          string
	Secret      string
	RedirectURL string
}

type google struct{ *client }

func newGoogle(app *GoogleApp, s *signer.Sessions) *google {
	c := newClient(
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

func (g *google) callback(req *http.Request) (string, error) {
	tok, err := g.client.token(req)
	if err != nil {
		return "", err
	}

	const url = "https://www.googleapis.com/oauth2/v3/userinfo"
	bs, err := g.client.get(tok, url)
	if err != nil {
		return "", err
	}

	var user struct {
		Email string `json:"email"`
	}
	if err := json.Unmarshal(bs, &user); err != nil {
		return "", err
	}
	ret := user.Email
	if ret == "" {
		return "", fmt.Errorf("empty login")
	}
	return ret, nil
}
