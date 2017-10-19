package aries

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

// C provides the request context for a web application.
type C struct {
	Path string
	User string

	App     string
	AppPath string

	Req  *http.Request
	Resp http.ResponseWriter

	HTTPS bool

	Data map[string]interface{}
}

// Redirect redirects the request to another URL.
func (c *C) Redirect(url string) {
	http.Redirect(c.Resp, c.Req, url, http.StatusFound)
}

// RespondJSON respond the request with a JSON object.
func (c *C) RespondJSON(dat interface{}) {
	enc := json.NewEncoder(c.Resp)
	if err := enc.Encode(dat); err != nil {
		log.Println(err)
	}
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

// ErrorStr returns an error to the request.
func (c *C) ErrorStr(code int, s string) {
	c.Error(code, errors.New(s))
}

// Errorf creates a formatted error and sends it to the client
// as a 500 error message.
func (c *C) Errorf(code int, f string, args ...interface{}) {
	c.Error(code, fmt.Errorf(f, args...))
}

// Error responds with a 500 error with the error message if err is not
// nil.
func (c *C) Error(code int, err error) bool {
	if err == nil {
		return false
	}
	http.Error(c.Resp, err.Error(), code)
	return true
}

// AltError responds with an error with alternative error meesage if err is not
// nil.
func (c *C) AltError(err error, code int, s string) bool {
	if err == nil {
		return false
	}
	AltError(c.Resp, err, s, code)
	return true
}

// IsMobile checks if the user agent of the request is mobile or not.
func (c *C) IsMobile() bool {
	return isMobile(c.Req.UserAgent())
}
