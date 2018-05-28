package oauth

import (
	"encoding/json"
	"fmt"

	"golang.org/x/oauth2"
	gh "golang.org/x/oauth2/github"
	"shanhu.io/aries"
	"shanhu.io/misc/signer"
)

// GitHubApp is the configuration of a GitHub Oauth App.
type GitHubApp struct {
	ID     string
	Secret string
}

type github struct{ client *Client }

func newGitHub(app *GitHubApp, s *signer.Sessions) *github {
	c := NewClient(
		&oauth2.Config{
			ClientID:     app.ID,
			ClientSecret: app.Secret,
			Scopes:       []string{}, // only need public information
			Endpoint:     gh.Endpoint,
		}, s,
	)
	return &github{client: c}
}

func (g *github) callback(c *aries.C) (string, error) {
	tok, err := g.client.Token(c)
	if err != nil {
		return "", err
	}

	bs, err := g.client.Get(c.Context, tok, "https://api.github.com/user")
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
