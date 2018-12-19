package aries

import (
	"testing"

	"net/http/httptest"

	"shanhu.io/base/httputil"
)

func TestStaticFiles(t *testing.T) {
	static := NewStaticFiles("testdata/static")

	s := httptest.NewServer(Serve(static))
	defer s.Close()

	c := httputil.NewClient(s.URL)
	for _, test := range []struct {
		p, want string
	}{
		{"/f1.html", "hello\n"},
		{"/f2.html", "hi\n"},
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
