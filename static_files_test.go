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
		p, want string
	}{
		{"/f1.html", "hello"},
		{"/f2.html", "hi"},
	} {
		reply, err := c.GetString(test.p)
		if err != nil {
			t.Errorf("%q - got error: %s", test.p, err)
			continue
		}
		if reply != test.want {
			t.Errorf("%q - want %q, got %q", test.p, test.want, reply)
		}
	}
}
