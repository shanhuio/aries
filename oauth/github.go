package oauth

import (
	"encoding/json"
	"fmt"
	"strconv"

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

type github struct{ c *Client }

func newGitHub(app *GitHubApp, s *signer.Sessions) *github {
	c := NewClient(
		&oauth2.Config{
			ClientID:     app.ID,
			ClientSecret: app.Secret,
			Scopes:       []string{}, // only need public information
			Endpoint:     gh.Endpoint,
		}, s,
	)
	return &github{c: c}
}

func (g *github) client() *Client { return g.c }

func (g *github) callback(c *aries.C) (*userMeta, *State, error) {
	tok, state, err := g.c.TokenState(c)
	if err != nil {
		return nil, nil, err
	}

	bs, err := g.c.Get(c.Context, tok, "https://api.github.com/user")
	if err != nil {
		return nil, nil, err
	}

	var user struct {
		ID int `json:"id"`
	}
	if err := json.Unmarshal(bs, &user); err != nil {
		return nil, nil, err
	}
	if user.ID == 0 {
		return nil, nil, fmt.Errorf("empty login")
	}
	meta := &userMeta{
		id: strconv.Itoa(user.ID),
	}
	return meta, state, nil
}
