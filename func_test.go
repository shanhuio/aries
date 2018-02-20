package aries

import (
	"testing"

	"fmt"
	"net/http/httptest"

	"smallrepo.com/base/httputil"
)

func TestFunc(t *testing.T) {
	const msg = "hello"
	f := func(c *C) error {
		fmt.Fprint(c.Resp, msg)
		return nil
	}
	s := httptest.NewServer(Func(f))
	defer s.Close()

	got, err := httputil.GetString(s.Client(), s.URL)
	if err != nil {
		t.Error(err)
		return
	}
	if got != msg {
		t.Errorf("want %q in response, got %s", msg, got)
	}
}

func TestFuncHTTPS(t *testing.T) {
	const msg = "hello"
	f := func(c *C) error {
		fmt.Fprint(c.Resp, msg)
		return nil
	}
	s := httptest.NewTLSServer(Func(f))
	defer s.Close()

	got, err := httputil.GetString(s.Client(), s.URL)
	if err != nil {
		t.Error(err)
		return
	}
	if got != msg {
		t.Errorf("want %q in response, got %s", msg, got)
	}
}
