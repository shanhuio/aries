package aries

import (
	"testing"

	"os"
	"io/ioutil"
	"path/filepath"
	"net/http/httptest"

	"smallrepo.com/base/httputil"
)

func TestTemplates(t *testing.T) {
	dir, err := ioutil.TempDir("", "aries")
	if err != nil {
		t.Fatal(err)
	}

	defer os.RemoveAll(dir)

	addFile := func(name, content string) {
		p := filepath.Join(dir, name)
		if err := ioutil.WriteFile(p, []byte(content), 0600); err != nil {
			t.Fatal(err)
		}
	}

	addFile("t1.html", "{{.Message1}}")
	addFile("t2.html", "{{.Message2}}")

	tmpls := NewTemplates(dir)

	f := func(c *C) error {
		dat := struct {
			Message1, Message2 string
		} {
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
	} {
		{"/t1.html", "hello"},
		{"/t2.html", "hi"},
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
