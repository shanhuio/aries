package aries

import (
	"testing"

	"fmt"
	"net/http/httptest"

	"shanhu.io/base/httputil"
)

func makeEchoRel(s string) Func {
	return func(c *C) error {
		fmt.Fprintf(c.Resp, "%s: %s", s, c.Rel())
		return nil
	}
}

func TestRouter(t *testing.T) {
	r := NewRouter()
	r.File("something", StringFunc("xxx"))
	r.Dir("books", makeEchoRel("books"))

	s := httptest.NewServer(Serve(r))
	defer s.Close()

	c := s.Client()
	host := s.URL

	for _, test := range []struct {
		p, want string
	}{
		{"/something", "xxx"},
		{"/books/xxx", "books: xxx"},
		{"/books/yyy", "books: yyy"},
		{"/books/", "books: "},
		{"/books", "books: "},
	} {
		got, err := httputil.GetString(c, host+test.p)
		if err != nil {
			t.Errorf("get %q, got error: %s", test.p, err)
			continue
		}

		if got != test.want {
			t.Errorf(
				"get %q, want %q in response, got %q",
				test.p, test.want, got,
			)
		}
	}

	for _, p := range []string{
		"/something/xxx",
		"/bookss",
		"/something/",
	} {
		code, err := httputil.GetCode(c, host+p)
		if err != nil {
			t.Error(err)
			continue
		}

		if code != 404 {
			t.Errorf("get %q, want 404 response, got %d", p, code)
		}
	}
}

func TestRouterWithIndex(t *testing.T) {
	r := NewRouter()
	r.Index(StringFunc("index"))

	sub := NewRouter()
	sub.Index(StringFunc("sub-index"))
	r.Dir("sub", sub.Serve)

	s := httptest.NewServer(Serve(r))
	defer s.Close()

	for _, test := range []struct {
		p, want string
	}{
		{"", "index"},
		{"/", "index"},
		{"/sub", "sub-index"},
		{"/sub/", "sub-index"},
	} {
		got, err := httputil.GetString(s.Client(), s.URL+test.p)
		if err != nil {
			t.Error(err)
			return
		}
		if got != test.want {
			t.Errorf(
				"get index page, want %q in response, got %q",
				test.want, got,
			)
		}
	}
}

func TestRouterWithDefault(t *testing.T) {
	r := NewRouter()
	r.Index(StringFunc("index"))
	r.Default(StringFunc("default"))
	s := httptest.NewServer(Serve(r))
	defer s.Close()

	got, err := httputil.GetString(s.Client(), s.URL+"/notfound")
	if err != nil {
		t.Error(err)
		return
	}

	if got != "default" {
		t.Errorf(
			"get a 404 page, want %q in response, got %q",
			"default", got,
		)
	}
}
