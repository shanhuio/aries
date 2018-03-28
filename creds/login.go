package creds

import (
	"fmt"
	"os"
	"time"

	"shanhu.io/aries"
	"shanhu.io/aries/oauth"
	"shanhu.io/misc/signer"
	"smallrepo.com/base/httputil"
)

// LoginWithKey uses the given PEM file to login a server, and returns the creds
// if succeess.
func LoginWithKey(p *EndPoint) (*Creds, error) {
	tty := !p.NoTTY

	k, err := readPrivateKey(p.PemFile, !p.NoPermCheck, tty)
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

	c := httputil.NewClient(p.Server)
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
	endPoint  *EndPoint
	credsFile string
	creds     *Creds // cached creds
}

// NewServerLogin returns a new server login with default user and pem file.
func NewServerLogin(s string) (*Login, error) {
	p, err := NewEndPoint(s)
	if err != nil {
		return nil, err
	}

	return NewLogin(p), nil
}

// NewLogin creates a new login stub with the given config.
func NewLogin(p *EndPoint) *Login {
	if p.User == "" {
		panic("user is empty")
	}

	cp := *p
	if cp.PemFile == "" {
		cp.PemFile = "key.pem"
	}

	return &Login{
		endPoint:  &cp,
		credsFile: Filename(p.Server) + ".json",
	}
}

// NewRobotLogin is a shorthand for NewLogin(NewRobot())
func NewRobotLogin(user, server, key string, env *aries.Env) *Login {
	return NewLogin(NewRobot(user, server, key, env))
}

func (lg *Login) readCreds() (*Creds, error) {
	if lg.endPoint.Homeless {
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
	if lg.endPoint.Homeless {
		panic("login server is homeless")
	}
	return WriteHomeJSONFile(lg.credsFile, cs)
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
	cs := lg.creds
	if cs == nil {
		if lg.endPoint.Homeless {
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
	pemFile := lg.endPoint.PemFile

	if !lg.endPoint.Homeless {
		var err error
		lg.endPoint.PemFile, err = HomeFile(pemFile)
		if err != nil {
			return nil, err
		}
	}

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
	if !lg.endPoint.Homeless {
		if err := lg.writeCreds(cs); err != nil {
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

	c := httputil.NewTokenClient(lg.endPoint.Server, tok)
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
	return httputil.NewTokenClient(server, tok), nil
}

// DialEndPoint creates a token client with the given endpoint.
func DialEndPoint(p *EndPoint) (*httputil.Client, error) {
	login := NewLogin(p)
	return login.Dial()
}
