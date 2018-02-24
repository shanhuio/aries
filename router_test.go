package aries

import (
	"testing"

	"fmt"
	"net/http/httptest"

	"smallrepo.com/base/httputil"
)

func makeEchoRel(s string) Func {
	return func(c *C) error {
		fmt.Fprintf(c.Resp, "%s: %s", s, c.Rel())
		return nil
	}
}

func TestRouter(t *testing.T) {
	r := NewRouter()
	r.File("something", MakeStringFunc("xxx"))
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
			t.Error(err)
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
	r.Index(MakeStringFunc("index"))
	s := httptest.NewServer(Serve(r))
	defer s.Close()

	got, err := httputil.GetString(s.Client(), s.URL)
	if err != nil {
		t.Error(err)
		return
	}
	if got != "index" {
		t.Errorf(
			"get index page, want %q in response, got %q",
			"index", got,
		)
	}
}

func TestRouterWithDefault(t *testing.T) {
	r := NewRouter()
	r.Index(MakeStringFunc("index"))
	r.Default(MakeStringFunc("default"))
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
