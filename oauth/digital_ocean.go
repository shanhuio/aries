package oauth

import (
	"github.com/digitalocean/godo"
	"golang.org/x/oauth2"
	"shanhu.io/aries"
	"shanhu.io/misc/signer"
)

// DigitalOceanApp stores the configuration of a Digital Ocean app.
type DigitalOceanApp struct {
	ID     string
	Secret string
}

type digitalOcean struct{ client *Client }

var digitalOceanEndpoint = oauth2.Endpoint{
	AuthURL:  "https://cloud.digitalocean.com/v1/oauth/authorize",
	TokenURL: "https://cloud.digitalocean.com/v1/oauth/token",
}

func newDigitalOcean(
	app *DigitalOceanApp, s *signer.Sessions,
) *digitalOcean {
	c := NewClient(
		&oauth2.Config{
			ClientID:     app.ID,
			ClientSecret: app.Secret,
			Endpoint:     digitalOceanEndpoint,
		}, s,
	)
	return &digitalOcean{client: c}
}

func (d *digitalOcean) callback(c *aries.C) (string, *State, error) {
	tok, state, err := d.client.TokenState(c)
	if err != nil {
		return "", nil, err
	}

	oc := oauth2.NewClient(c.Context, oauth2.StaticTokenSource(tok))
	client := godo.NewClient(oc)

	account, _, err := client.Account.Get(c.Context)
	if err != nil {
		return "", nil, err
	}

	return account.Email, state, nil
}
