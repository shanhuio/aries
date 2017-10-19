package oauth

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"golang.org/x/oauth2"

	"shanhu.io/misc/signer"
)

type client struct {
	config *oauth2.Config
	states *signer.Sessions
}

func newClient(c *oauth2.Config, states *signer.Sessions) *client {
	return &client{
		config: c,
		states: states,
	}
}

func (c *client) signInURL() string {
	return c.config.AuthCodeURL(c.states.NewState())
}

func (c *client) token(req *http.Request) (*oauth2.Token, error) {
	state, code := stateCode(req)
	if state == "" {
		return nil, fmt.Errorf("invalid oauth redirect")
	}

	check := c.states.CheckState(state)
	if !check {
		return nil, fmt.Errorf("state invalid")
	}

	tok, err := c.config.Exchange(oauth2.NoContext, code)
	if err != nil {
		return nil, fmt.Errorf("exchange failed: %v", err)
	}
	if !tok.Valid() {
		return nil, fmt.Errorf("token is invalid")
	}
	return tok, nil
}

func (c *client) get(tok *oauth2.Token, url string) ([]byte, error) {
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
