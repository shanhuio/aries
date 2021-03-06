package aries

import (
	"testing"

	"shanhu.io/aries/https/httpstest"
	"shanhu.io/misc/httputil"
)

func TestHostMux(t *testing.T) {
	m := NewHostMux()
	m.Set("shanhu.io", StringFunc("shanhu"))
	m.Set("h8liu.io", StringFunc("h8liu"))

	s, err := httpstest.NewServer([]string{
		"shanhu.io", "h8liu.io",
	}, Serve(m))
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	c := s.Client()

	for _, test := range []struct {
		url, want string
	}{
		{"https://shanhu.io", "shanhu"},
		{"https://h8liu.io", "h8liu"},
	} {
		got, err := httputil.GetString(c, test.url)
		if err != nil {
			t.Errorf("get %q, got error %s", test.url, err)
		} else if got != test.want {
			t.Errorf("get %q, got %q, want %q", test.url, got, test.want)
		}
	}
}
