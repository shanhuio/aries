package aries

import (
	"testing"

	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
)

func httpGetString(c *http.Client, url string) (string, error) {
	resp, err := c.Get(url)
	if err != nil {
		return "", err
	}
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(bs), nil
}

func TestHandler(t *testing.T) {
	msg := "hello"
	f := func(c *C) { fmt.Fprint(c.Resp, msg) }
	s := httptest.NewServer(HandlerFunc(f, false))
	defer s.Close()

	got, err := httpGetString(s.Client(), s.URL)
	if err != nil {
		t.Error(err)
		return
	}
	if got != msg {
		t.Errorf("want %q in response, got %s", msg, got)
	}
}

func TestHTTPSHandler(t *testing.T) {
	msg := "hello"
	f := func(c *C) { fmt.Fprint(c.Resp, msg) }
	s := httptest.NewTLSServer(HandlerFunc(f, false))
	defer s.Close()

	got, err := httpGetString(s.Client(), s.URL)
	if err != nil {
		t.Error(err)
		return
	}
	if got != msg {
		t.Errorf("want %q in response, got %s", msg, got)
	}
}
