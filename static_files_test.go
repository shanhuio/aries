package aries

import (
	"testing"

	"io/ioutil"
	"net/http/httptest"
	"os"
	"path/filepath"

	"smallrepo.com/base/httputil"
)

func TestStaticFiles(t *testing.T) {
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

	addFile("f1.html", "hello")
	addFile("f2.html", "hi")
	static := NewStaticFiles(dir)

	s := httptest.NewServer(Serve(static))
	defer s.Close()

	c := httputil.NewClient(s.URL)
	for _, test := range []struct {
		url, want string
	}{
		{"/f1.html", "hello"},
		{"/f2.html", "hi"},
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
