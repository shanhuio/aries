package creds

import (
	"crypto/rsa"
	"fmt"
	"time"

	"shanhu.io/aries"
	"shanhu.io/aries/oauth"
	"shanhu.io/misc/errcode"
	"shanhu.io/misc/httputil"
	"shanhu.io/misc/signer"
)

func readEndpointKey(p *EndPoint) (*rsa.PrivateKey, error) {
	tty := !p.NoTTY
	if p.Key != nil {
		return ParsePrivateKey("key", p.Key, tty)
	}
	return readPrivateKey(p.PemFile, !p.NoPermCheck, tty)
}

// LoginWithKey uses the given PEM file to login a server, and returns the creds
// if succeess.
func LoginWithKey(p *EndPoint) (*Creds, error) {
	k, err := readEndpointKey(p)
	if err != nil {
		return nil, err
	}

	signed, err := signer.RSASignTime(k)
	if err != nil {
		return nil, err
	}

	req := &oauth.LoginRequest{
		User:       p.User,
		SignedTime: signed,
	}
	cs := &Creds{Server: p.Server}

	c, err := httputil.NewClient(p.Server)
	if err != nil {
		return nil, err
	}
	if p.Transport != nil {
		c.Transport = p.Transport
	}
	if err := c.JSONCall("/pubkey/signin", req, &cs.Creds); err != nil {
		return nil, err
	}

	if cs.Creds.User != p.User {
		return nil, fmt.Errorf("login as user %q, got %q", p.User, cs.User)
	}

	return cs, nil
}

// Login is a helper stub to perform login actions.
type Login struct {
	endPoint   *EndPoint
	credsStore credsStore
	creds      *Creds // cached creds
}

// NewServerLogin returns a new server login with default user and pem file.
func NewServerLogin(s string) (*Login, error) {
	p, err := NewEndPoint(s)
	if err != nil {
		return nil, err
	}
	return NewLogin(p)
}

// NewLogin creates a new login stub with the given config.
func NewLogin(p *EndPoint) (*Login, error) {
	if p.User == "" {
		return nil, errcode.InvalidArgf("user is empty")
	}

	cp := *p
	if cp.PemFile == "" && !cp.Homeless {
		pem, err := HomeFile("key.pem")
		if err != nil {
			return nil, errcode.Internalf("fail to get home: %v", err)
		}
		cp.PemFile = pem
	}

	lg := &Login{endPoint: &cp}
	if !p.Homeless {
		lg.credsStore = newHomeCredsStore(p.Server)
	}
	return lg, nil
}

// NewRobotLogin is a shorthand for NewLogin(NewRobot())
func NewRobotLogin(
	user, server, key string, env *aries.Env,
) (*Login, error) {
	return NewLogin(NewRobot(user, server, key, env))
}

func (lg *Login) check(cs *Creds) (bool, error) {
	if cs.User != lg.endPoint.User {
		return false, nil
	}
	if cs.Server != lg.endPoint.Server {
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
	if lg.endPoint.Homeless {
		// Nothing cached anywhere, just return a new one.
		return lg.GetToken()
	}

	cs := lg.creds
	if cs == nil && lg.credsStore != nil {
		newCreds, err := lg.credsStore.read()
		if err != nil {
			if errcode.IsNotFound(err) {
				return lg.GetToken()
			}
			return "", err
		}
		if newCreds == nil {
			panic("should have creds loaded from the file system")
		}
		cs = newCreds
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
	return LoginWithKey(lg.endPoint)
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
	if lg.credsStore != nil {
		if err := lg.credsStore.write(cs); err != nil {
			return "", err
		}
	}
	return cs.Creds.Token, nil
}

// Dial creates an token client.
func (lg *Login) Dial() (*httputil.Client, error) {
	tok, err := lg.Token()
	if err != nil {
		return nil, err
	}

	c, err := httputil.NewTokenClient(lg.endPoint.Server, tok)
	if err != nil {
		return nil, err
	}
	c.Transport = lg.endPoint.Transport
	return c, nil
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
	return httputil.NewTokenClient(server, tok)
}

// DialEndPoint creates a token client with the given endpoint.
func DialEndPoint(p *EndPoint) (*httputil.Client, error) {
	login, err := NewLogin(p)
	if err != nil {
		return nil, err
	}
	return login.Dial()
}
