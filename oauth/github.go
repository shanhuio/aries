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
	ID          string
	Secret      string
	RedirectURL string
}

type github struct {
	c          *Client
	queryEmail bool
}

const githubEmailScope = "user:email"

func newGitHubWithScopes(
	app *GitHubApp, s *signer.Sessions, scopes []string,
) *github {
	if scopes == nil {
		scopes = []string{}
	}
	queryEmail := false
	for _, scope := range scopes {
		if scope == githubEmailScope {
			queryEmail = true
		}
	}
	c := NewClient(
		&oauth2.Config{
			ClientID:     app.ID,
			ClientSecret: app.Secret,
			Scopes:       scopes, // only need public information
			Endpoint:     gh.Endpoint,
			RedirectURL:  app.RedirectURL,
		}, s,
	)
	return &github{c: c, queryEmail: queryEmail}
}

func newGitHubWithEmail(app *GitHubApp, s *signer.Sessions) *github {
	return newGitHubWithScopes(app, s, []string{githubEmailScope})
}

func newGitHub(app *GitHubApp, s *signer.Sessions) *github {
	return newGitHubWithScopes(app, s, nil)
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
		Login string `json:"login"`
		ID    int    `json:"id"`
	}
	if err := json.Unmarshal(bs, &user); err != nil {
		return nil, nil, err
	}
	if user.ID == 0 {
		return nil, nil, fmt.Errorf("empty login")
	}

	var email string
	if g.queryEmail {
		const url = "https://api.github.com/user/emails"
		bs, err := g.c.Get(c.Context, tok, url)
		if err != nil {
			return nil, nil, err
		}

		type userEmail struct {
			Email    string `json:"email"`
			Verified bool   `json:"verified"`
			Primary  bool   `json:"primary"`
		}

		var emails []*userEmail
		if err := json.Unmarshal(bs, &emails); err != nil {
			return nil, nil, err
		}
		for _, m := range emails {
			if m.Primary && m.Verified {
				email = m.Email
			}
		}
	}
	meta := &userMeta{
		id:    strconv.Itoa(user.ID),
		name:  user.Login,
		email: email,
	}
	return meta, state, nil
}
