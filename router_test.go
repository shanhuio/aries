package aries

import (
	"testing"

	"fmt"
	"net/http/httptest"
)

func makeEcho(s string) Func {
	return func(c *C) error {
		fmt.Fprint(c.Resp, s)
		return nil
	}
}

func makeEchoRel(s string) Func {
	return func(c *C) error {
		fmt.Fprintf(c.Resp, "%s: %s", s, c.Rel())
		return nil
	}
}

func TestRouter(t *testing.T) {
	r := NewRouter()
	r.File("something", makeEcho("xxx"))
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
		got, err := httpGetString(c, host+test.p)
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
		code, err := httpGetCode(c, host+p)
		if err != nil {
			t.Error(err)
			continue
		}

		if code != 404 {
			t.Errorf("get %q, want 404 response, got %d", p, code)
		}
	}

}
