// Package mailman provides an HTTP Oauth2 based module that sends email using
// Gmail API.
package mailman

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
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
type Tokens interface {
	Get(ctx context.Context, email string) (*oauth2.Token, error)
	Set(ctx context.Context, email string, t *oauth2.Token) error
}

// Mailman is a http server module for sending emails using gmail's
// OAuth2 API.
type Mailman struct {
	config *oauth2.Config
	client *oauth.Client
	tokens Tokens
}

// Config contains configuration for a mailman.
type Config struct {
	App      *oauth.GoogleApp
	StateKey []byte
	Tokens   Tokens
}

// New creates a new mailman.
func New(c *Config) *Mailman {
	states := signer.NewSessions(c.StateKey, time.Minute*3)

	const gmailSendScope = "https://www.googleapis.com/auth/gmail.send"
	const emailScope = "https://www.googleapis.com/auth/userinfo.email"

	oc := &oauth2.Config{
		ClientID:     c.App.ID,
		ClientSecret: c.App.Secret,
		Scopes: []string{
			gmailSendScope,
			emailScope,
		},
		Endpoint:    goauth2.Endpoint,
		RedirectURL: c.App.RedirectURL,
	}

	return &Mailman{
		config: oc,
		client: oauth.NewClient(oc, states, oauth.MethodGoogle),
		tokens: c.Tokens,
	}
}

func (m *Mailman) signInURL() string {
	return m.client.OfflineSignInURL(new(oauth.State))
}

func (m *Mailman) tokenState(c *aries.C) (*oauth2.Token, *oauth.State, error) {
	return m.client.TokenState(c)
}

func (m *Mailman) serveIndex(c *aries.C) error {
	emails, ok := c.Req.URL.Query()["email"]
	if !ok {
		c.Resp.Write([]byte("Please specify email parameter"))
		return nil
	}

	// We simply take the first specified email parameter.
	tok, err := m.tokens.Get(c.Context, emails[0])
	if err != nil {
		if errcode.IsNotFound(err) {
			c.Resp.Write([]byte("Token not found"))
			return nil
		}
		return err
	}
	return aries.PrintJSON(c, tok)
}

// Send sends an email. Needs to setup OAuth2 first.
func (m *Mailman) Send(
	ctx context.Context, from string, body []byte) (string, error) {
	tok, err := m.tokens.Get(ctx, from)
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

	u := &url.URL{
		Scheme: "https",
		Host:   "www.googleapis.com",
	}
	client := &httputil.Client{Server: u, Token: curTok.AccessToken}

	const route = "/gmail/v1/users/me/messages/send?alt=json"
	if err := client.JSONCall(route, &msg, &resp); err != nil {
		return "", err
	}

	return resp.ID, nil
}

// SendRequest is an request for sending a mail.
type SendRequest struct {
	From string
	Body []byte
}

func (m *Mailman) apiSend(c *aries.C, req *SendRequest) (string, error) {
	return m.Send(c.Context, req.From, req.Body)
}

func (m *Mailman) serveCallback(c *aries.C) error {
	token, _, err := m.tokenState(c)
	if err != nil {
		return err
	}

	if token.RefreshToken == "" {
		return fmt.Errorf("refresh token empty")
	}

	// Let's find out which email this callback is for.
	const userInfoURL = "https://www.googleapis.com/oauth2/v3/userinfo"
	bs, err := m.client.Get(c.Context, token, userInfoURL)
	if err != nil {
		return err
	}

	var user struct {
		Email string `json:"email"`
	}
	if err := json.Unmarshal(bs, &user); err != nil {
		return err
	}

	if err := m.tokens.Set(c.Context, user.Email, token); err != nil {
		return err
	}

	// redirect to index with email parameter
	c.Redirect(path.Dir(c.Path) + "?email=" + user.Email)
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
