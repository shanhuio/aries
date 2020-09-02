package oauth

import (
	"encoding/json"
	"fmt"
	"strconv"

	"golang.org/x/oauth2"
	gh "golang.org/x/oauth2/github"
	"shanhu.io/aries"
	"shanhu.io/misc/signer"
	"shanhu.io/misc/strutil"
)

// GitHubApp is the configuration of a GitHub Oauth App.
type GitHubApp struct {
	ID          string
	Secret      string
	RedirectURL string

	WithEmail bool
	Scopes    []string
}

type github struct {
	c          *Client
	queryEmail bool
}

const githubEmailScope = "user:email"

func newGitHub(app *GitHubApp, s *signer.Sessions) *github {
	scopeSet := make(map[string]bool)
	if app.WithEmail {
		scopeSet[githubEmailScope] = true
	}
	for _, scope := range app.Scopes {
		scopeSet[scope] = true
	}
	scopes := strutil.SortedList(scopeSet)
	if scopes == nil {
		scopes = []string{}
	}
	queryEmail := scopeSet[githubEmailScope]

	c := NewClient(
		&oauth2.Config{
			ClientID:     app.ID,
			ClientSecret: app.Secret,
			Scopes:       scopes, // only need public information
			Endpoint:     gh.Endpoint,
			RedirectURL:  app.RedirectURL,
		}, s, MethodGitHub,
	)
	return &github{c: c, queryEmail: queryEmail}
}

func (g *github) client() *Client { return g.c }

func (g *github) callback(c *aries.C) (*UserMeta, *State, error) {
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
	meta := &UserMeta{
		Method: MethodGitHub,
		ID:     strconv.Itoa(user.ID),
		Name:   user.Login,
		Email:  email,
	}
	return meta, state, nil
}
