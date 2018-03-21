package aries

import (
	"testing"

	"net/http/httptest"

	"smallrepo.com/base/httputil"
)

func TestTemplates(t *testing.T) {
	tmpls := NewTemplates("testdata/templates", nil)

	f := func(c *C) error {
		dat := struct {
			Message1, Message2 string
		}{
			Message1: "hello",
			Message2: "hi",
		}
		return tmpls.Serve(c, c.Rel(), dat)
	}

	s := httptest.NewServer(Func(f))
	defer s.Close()

	c := httputil.NewClient(s.URL)
	for _, test := range []struct {
		url, want string
	}{
		{"/t1.html", "hello\n"},
		{"/t2.html", "hi\n"},
	} {
		reply, err := c.GetString(test.url)
		if err != nil {
			t.Errorf(
				"http get %q, got error: %s",
				test.url, err,
			)
		}
		if reply != test.want {
			t.Errorf(
				"http get %q, want %q, got %q",
				test.url, test.want, reply,
			)
		}
	}
}

func TestTemplatesJSON(t *testing.T) {
	tmpls := NewTemplates(TemplatesJSON, nil)

	f := func(c *C) error {
		dat := struct {
			Message1, Message2 string
		}{
			Message1: "hello",
			Message2: "hi",
		}
		return tmpls.Serve(c, c.Rel(), dat)
	}

	s := httptest.NewServer(Func(f))
	defer s.Close()

	c := httputil.NewClient(s.URL)
	test := struct {
		url, want string
	}{
		"/", `{"Message1":"hello","Message2":"hi"}`,
	}

	reply, err := c.GetString(test.url)
	if err != nil {
		t.Errorf(
			"http get %q, got error: %s",
			test.url, err,
		)
	}
	if reply != test.want {
		t.Errorf(
			"http get %q, want %q, got %q",
			test.url, test.want, reply,
		)
	}
}
