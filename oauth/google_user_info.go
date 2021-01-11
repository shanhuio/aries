package oauth

import (
	"context"
	"encoding/json"

	"golang.org/x/oauth2"
)

// GoogleUserInfo stores a Google user's basic personal info.
type GoogleUserInfo struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

// GetGoogleUserInfo queries Google OAuth endpoint for user info data.
func GetGoogleUserInfo(
	ctx context.Context, c *Client, tok *oauth2.Token,
) (*GoogleUserInfo, error) {
	bs, err := c.Get(ctx, tok, "https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		return nil, err
	}

	user := new(GoogleUserInfo)
	if err := json.Unmarshal(bs, user); err != nil {
		return nil, err
	}

	return user, nil
}
