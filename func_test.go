package aries

import (
	"testing"

	"net/http/httptest"

	"shanhu.io/misc/httputil"
)

func TestFunc(t *testing.T) {
	const msg = "hello"
	f := StringFunc(msg)
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
	f := StringFunc(msg)
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
