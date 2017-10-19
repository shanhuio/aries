package oauth

import (
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
	gh "golang.org/x/oauth2/github"

	"shanhu.io/misc/signer"
)

// GitHubApp is the configuration of a GitHub Oauth App.
type GitHubApp struct {
	ID     string
	Secret string
}

type github struct{ *client }

func newGitHub(app *GitHubApp, s *signer.Sessions) *github {
	c := newClient(
		&oauth2.Config{
			ClientID:     app.ID,
			ClientSecret: app.Secret,
			Scopes:       []string{}, // only need public information
			Endpoint:     gh.Endpoint,
		}, s,
	)
	return &github{c}
}

func (g *github) callback(req *http.Request) (string, error) {
	tok, err := g.client.token(req)
	if err != nil {
		return "", err
	}

	bs, err := g.client.get(tok, "https://api.github.com/user")
	if err != nil {
		return "", err
	}

	var user struct {
		Login string `json:"login"`
	}
	if err := json.Unmarshal(bs, &user); err != nil {
		return "", err
	}
	ret := user.Login
	if ret == "" {
		return "", fmt.Errorf("empty login")
	}
	return ret, nil
}
