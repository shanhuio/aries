package oauth

import (
	"context"
	"fmt"
	"io/ioutil"

	"golang.org/x/oauth2"
	"shanhu.io/aries"
	"shanhu.io/misc/signer"
)

// State contains a JSON marshalable state for OAuth2 sign in.
type State struct {
	Dest string // URL to redirect after signing in.
}

// Client is an oauth client for oauth2 exchanges.
type Client struct {
	config *oauth2.Config
	states *signer.Sessions
}

// NewClient creates a new oauth client for oauth2 exchnages.
func NewClient(c *oauth2.Config, states *signer.Sessions) *Client {
	return &Client{
		config: c,
		states: states,
	}
}

// SignInURL returns the online signin URL for redirection.
func (c *Client) SignInURL(s *State) string {
	state, _, err := c.states.NewJSON(s)
	if err != nil {
		panic(err)
	}
	return c.config.AuthCodeURL(state)
}

// OfflineSignInURL returns the offline signin URL for redirection.
func (c *Client) OfflineSignInURL(s *State) string {
	state, _, err := c.states.NewJSON(s)
	if err != nil {
		panic(err)
	}
	return c.config.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

// TokenState extracts the oauth2 access token and state from the request.
func (c *Client) TokenState(ctx *aries.C) (*oauth2.Token, *State, error) {
	stateStr, code := stateCode(ctx.Req)
	if stateStr == "" {
		return nil, nil, fmt.Errorf("invalid oauth redirect")
	}

	state := new(State)
	if !c.states.CheckJSON(stateStr, state) {
		return nil, nil, fmt.Errorf("state invalid")
	}

	tok, err := c.config.Exchange(ctx.Context, code)
	if err != nil {
		return nil, nil, fmt.Errorf("exchange failed: %v", err)
	}
	if !tok.Valid() {
		return nil, nil, fmt.Errorf("token is invalid")
	}
	return tok, state, nil
}

// Get gets an URL using the given token.
func (c *Client) Get(
	ctx context.Context, tok *oauth2.Token, url string,
) ([]byte, error) {
	callClient := c.config.Client(oauth2.NoContext, tok)
	resp, err := callClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("oauth2 get: %v", err)
	}
	defer resp.Body.Close()

	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("oauth2 read body: %v", err)
	}
	return bs, nil
}
