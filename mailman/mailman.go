// Package mailman provides an HTTP Oauth2 based module that sends email using
// Gmail API.
package mailman

import (
	"context"
	"encoding/base64"
	"fmt"
	"path"
	"time"

	"golang.org/x/oauth2"
	goauth2 "golang.org/x/oauth2/google"
	"shanhu.io/aries"
	"shanhu.io/aries/oauth"
	"shanhu.io/misc/errcode"
	"shanhu.io/misc/httputil"
	"shanhu.io/misc/signer"
)

// Token is an interface that gets and fetches an OAuth2 refresh token.
type Token interface {
	Get(ctx context.Context) (*oauth2.Token, error)
	Set(ctx context.Context, t *oauth2.Token) error
}

// Mailman is a http server module for sending emails using gmail's
// OAuth2 API.
type Mailman struct {
	config *oauth2.Config
	client *oauth.Client
	token  Token
}

// Config contains configuration for a mailman.
type Config struct {
	App      *oauth.GoogleApp
	StateKey []byte
	Token    Token
}

// New creates a new mailman.
func New(c *Config) *Mailman {
	states := signer.NewSessions(c.StateKey, time.Minute*3)

	const gmailSendScope = "https://www.googleapis.com/auth/gmail.send"

	oc := &oauth2.Config{
		ClientID:     c.App.ID,
		ClientSecret: c.App.Secret,
		Scopes:       []string{gmailSendScope},
		Endpoint:     goauth2.Endpoint,
		RedirectURL:  c.App.RedirectURL,
	}

	return &Mailman{
		config: oc,
		client: oauth.NewClient(oc, states),
		token:  c.Token,
	}
}

func (m *Mailman) signInURL() string {
	return m.client.OfflineSignInURL(new(oauth.State))
}

func (m *Mailman) tokenState(c *aries.C) (*oauth2.Token, *oauth.State, error) {
	return m.client.TokenState(c)
}

func (m *Mailman) serveIndex(c *aries.C) error {
	tok, err := m.token.Get(c.Context)
	if err != nil {
		if errcode.IsNotFound(err) {
			return fmt.Errorf("mailman token not found")
		}
		return err
	}
	return aries.PrintJSON(c, tok)
}

// Send sends an email. Needs to setup OAuth2 first.
func (m *Mailman) Send(ctx context.Context, body []byte) (string, error) {
	tok, err := m.token.Get(ctx)
	if err != nil {
		if errcode.IsNotFound(err) {
			return "", fmt.Errorf("mailman not setup yet")
		}
		return "", err
	}

	// refresh the token.
	curTok, err := m.config.TokenSource(ctx, tok).Token()
	if err != nil {
		return "", err
	}

	var msg struct {
		Raw string `json:"raw"`
	}
	msg.Raw = base64.URLEncoding.EncodeToString(body)

	var resp struct {
		ID string `json:"id"`
	}

	const url = "https://www.googleapis.com/"
	client := httputil.NewTokenClient(url, curTok.AccessToken)

	const route = "/gmail/v1/users/me/messages/send?alt=json"
	if err := client.JSONCall(route, &msg, &resp); err != nil {
		return "", err
	}

	return resp.ID, nil
}

// SendRequest is an request for sending a mail.
type SendRequest struct {
	Body []byte
}

func (m *Mailman) apiSend(c *aries.C, req *SendRequest) (string, error) {
	return m.Send(c.Context, req.Body)
}

func (m *Mailman) serveCallback(c *aries.C) error {
	token, _, err := m.tokenState(c)
	if err != nil {
		return err
	}

	if token.RefreshToken == "" {
		return fmt.Errorf("refresh token empty")
	}
	if err := m.token.Set(c.Context, token); err != nil {
		return err
	}

	c.Redirect(path.Dir(c.Path)) // redirect to index
	return nil
}

func (m *Mailman) serveSetup(c *aries.C) error {
	c.Redirect(m.signInURL())
	return nil
}

// Router returns the mailman module router.
func (m *Mailman) Router() *aries.Router {
	r := aries.NewRouter()
	r.Index(m.serveIndex)
	r.JSONCallMust("send", m.apiSend)
	r.File("callback", m.serveCallback)
	r.File("setup", m.serveSetup)
	return r
}
