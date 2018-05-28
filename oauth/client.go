package oauth

import (
	"context"
	"fmt"
	"io/ioutil"

	"golang.org/x/oauth2"
	"shanhu.io/aries"
	"shanhu.io/misc/signer"
)

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
func (c *Client) SignInURL() string {
	return c.config.AuthCodeURL(c.states.NewState())
}

// OfflineSignInURL returns the offline signin URL for redirection.
func (c *Client) OfflineSignInURL() string {
	state := c.states.NewState()
	return c.config.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

// Token extracts the oauth2 access token from the request.
func (c *Client) Token(ctx *aries.C) (*oauth2.Token, error) {
	state, code := stateCode(ctx.Req)
	if state == "" {
		return nil, fmt.Errorf("invalid oauth redirect")
	}

	check := c.states.CheckState(state)
	if !check {
		return nil, fmt.Errorf("state invalid")
	}

	tok, err := c.config.Exchange(ctx.Context, code)
	if err != nil {
		return nil, fmt.Errorf("exchange failed: %v", err)
	}
	if !tok.Valid() {
		return nil, fmt.Errorf("token is invalid")
	}
	return tok, nil
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
