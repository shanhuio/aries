package creds

import (
	"fmt"
	"os"
	"time"

	"shanhu.io/aries/oauth"
	"shanhu.io/misc/signer"
	"smallrepo.com/base/httputil"
)

// LoginWithKey uses the given PEM file to login a server, and returns the creds
// if succeess.
func LoginWithKey(user, server, pemFile string, tty bool) (*Creds, error) {
	k, err := ReadPrivateKey(pemFile, tty)
	if err != nil {
		return nil, err
	}

	signed, err := signer.RSASignTime(k)
	if err != nil {
		return nil, err
	}

	req := &oauth.LoginRequest{
		User:       user,
		SignedTime: signed,
	}
	cs := &Creds{Server: server}

	c := httputil.NewClient(server)
	if err := c.JSONCall("pubkey/signin", req, &cs.Creds); err != nil {
		return nil, err
	}

	if cs.Creds.User != user {
		return nil, fmt.Errorf("login as user %q, got %q", user, cs.User)
	}

	return cs, nil
}

// Login is a helper stub to perform login actions.
type Login struct {
	server    string
	user      string
	pemFile   string
	credsFile string
	homeless  bool

	tty bool // if it is running under a terminal

	creds *Creds // cached creds
}

// NewServerLogin returns a new server login with default user and pem file.
func NewServerLogin(s string) (*Login, error) {
	u, err := currentUser()
	if err != nil {
		return nil, err
	}

	return NewLogin(&LoginConfig{User: u, Server: s}), nil
}

// NewLogin creates a new login stub with the given config.
func NewLogin(c *LoginConfig) *Login {
	if c.User == "" {
		panic("user is empty")
	}

	pemFile := c.PemFile
	if pemFile == "" {
		pemFile = "key.pem"
	}

	return &Login{
		server:    c.Server,
		user:      c.User,
		pemFile:   pemFile,
		credsFile: Filename(c.Server) + ".json",
		homeless:  c.Homeless,
		tty:       !c.NoTTY,
	}
}

func (lg *Login) readCreds() (*Creds, error) {
	if lg.homeless {
		panic("login server is homeless")
	}

	ret := &Creds{}
	if err := ReadHomeJSONFile(lg.credsFile, ret); err != nil {
		return nil, err
	}
	lg.creds = ret
	return ret, nil
}

func (lg *Login) writeCreds(cs *Creds) error {
	if lg.homeless {
		panic("login server is homeless")
	}
	return WriteHomeJSONFile(lg.credsFile, cs)
}

func (lg *Login) check(cs *Creds) (bool, error) {
	if cs.User != lg.user {
		return false, nil
	}
	if cs.Server != lg.server {
		return false, nil
	}

	expires := time.Unix(0, cs.Creds.Expires)
	now := time.Now()
	if !now.Before(expires) {
		return false, nil
	}

	return true, nil
}

// Token returns the login token for the login. If a valid token is already
// cached, it returns the cached one.
func (lg *Login) Token() (string, error) {
	cs := lg.creds
	if cs == nil {
		if lg.homeless {
			// Nothing cached anywhere, just return a new one.
			return lg.GetToken()
		}

		// Try read the cache on file system.
		var err error
		if cs, err = lg.readCreds(); err != nil {
			if os.IsNotExist(err) {
				return lg.GetToken()
			}
			return "", err
		}
		if cs == nil {
			panic("should have creds loaded from the file system")
		}
		lg.creds = cs
	}

	// now we loaded a cached creds
	ok, err := lg.check(cs)
	if err != nil {
		return "", err
	}
	if !ok {
		return lg.GetToken()
	}

	return cs.Token, nil
}

// Do performs the login and returns the credentials.
// It does not read or write the credential cache file.
func (lg *Login) Do() (*Creds, error) {
	pemFile := lg.pemFile

	if !lg.homeless {
		var err error
		pemFile, err = HomeFile(pemFile)
		if err != nil {
			return nil, err
		}
	}

	return LoginWithKey(lg.user, lg.server, pemFile, lg.tty)
}

// GetToken returns the login token for the login. It ignores and overwrites
// any existing login token that uses the same login creds file.
func (lg *Login) GetToken() (string, error) {
	cs, err := lg.Do()
	if err != nil {
		return "", err
	}

	// cache it
	lg.creds = cs

	// If not homeless, also cache it in home directory.
	if !lg.homeless {
		if err := lg.writeCreds(cs); err != nil {
			return "", err
		}
	}
	return cs.Creds.Token, nil
}

// LoginServer uses the default setting to login into a server.
func LoginServer(server string) (string, error) {
	login, err := NewServerLogin(server)
	if err != nil {
		return "", err
	}
	return login.Token()
}

// Dial logins the server and returns the httputil client.
func Dial(server string) (*httputil.Client, error) {
	tok, err := LoginServer(server)
	if err != nil {
		return nil, err
	}
	return httputil.NewTokenClient(server, tok), nil
}
