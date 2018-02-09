package aries

import (
	"net/http"
	"strings"
	"time"

	"shanhu.io/misc/errcode"
)

// C provides the request context for a web application.
type C struct {
	Path string

	User      string
	UserLevel int // 0 for normal user. 0 with empty User is anonymous.

	Req  *http.Request
	Resp http.ResponseWriter

	HTTPS bool

	Data map[string]interface{}

	route    *route
	routePos int
}

// NewContext creates a new context from the incomming request.
func NewContext(w http.ResponseWriter, req *http.Request) *C {
	isHTTPS := false
	u := req.URL
	if strings.ToLower(u.Scheme) == "https" {
		isHTTPS = true
	}
	if strings.ToLower(req.Header.Get("X-Forwared-Proto")) == "https" {
		isHTTPS = true
	}

	return &C{
		Path:  u.Path,
		Resp:  w,
		Req:   req,
		HTTPS: isHTTPS,
		Data:  make(map[string]interface{}),

		route: newRoute(u.Path),
	}
}

// Redirect redirects the request to another URL.
func (c *C) Redirect(url string) {
	http.Redirect(c.Resp, c.Req, url, http.StatusFound)
}

// Rel returns the current relative route. The return value changes if the
// routing is using a router, otherwise, it will always return the full routing
// path.
func (c *C) Rel() string {
	return c.route.rel(c.routePos)
}

// ReadCookie reads the cookie from the context.
func (c *C) ReadCookie(name string) string {
	cookie, err := c.Req.Cookie(name)
	if err != nil || cookie == nil {
		return ""
	}
	return cookie.Value
}

// WriteCookie sets a cookie.
func (c *C) WriteCookie(name, v string, expires time.Time) {
	cookie := &http.Cookie{
		Name:    name,
		Value:   v,
		Path:    "/",
		Expires: expires,
		Secure:  c.HTTPS,
	}
	http.SetCookie(c.Resp, cookie)
}

// ClearCookie clears a cookie.
func (c *C) ClearCookie(name string) {
	cookie := &http.Cookie{
		Name:   name,
		Value:  "",
		Path:   "/",
		Secure: c.HTTPS,
	}
	http.SetCookie(c.Resp, cookie)
}

// ErrCode returns an error based on its error code.
func (c *C) ErrCode(err error) bool {
	if err == nil {
		return false
	}
	code := errcode.Of(err)
	switch code {
	case errcode.NotFound:
		return c.replyError(404, err)
	case errcode.Internal:
		return c.replyError(500, err)
	case errcode.Unauthorized:
		return c.replyError(403, err)
	case errcode.InvalidArg:
		return c.replyError(400, err)
	}
	return c.replyError(500, err)
}

func (c *C) replyError(code int, err error) bool {
	if err == nil {
		return false
	}
	http.Error(c.Resp, err.Error(), code)
	return true
}

// IsMobile checks if the user agent of the request is mobile or not.
func (c *C) IsMobile() bool {
	return isMobile(c.Req.UserAgent())
}
