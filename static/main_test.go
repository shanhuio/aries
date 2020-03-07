package static

import (
	"testing"

	"shanhu.io/aries"
	"shanhu.io/aries/ariestest"
	"shanhu.io/misc/httputil"
)

func TestMain(t *testing.T) {
	config := &config{Dir: "testdata"}
	service, err := main(&aries.Env{Config: config})
	if err != nil {
		t.Fatal(err)
	}

	s, err := ariestest.HTTPSServer("shanhu.io", service)
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	c := httputil.NewClient("https://shanhu.io")
	c.Transport = s.Transport

	str, err := c.GetString("/")
	if err != nil {
		t.Fatal(err)
	}
	const want = "hello\n"
	if str != want {
		t.Errorf("get /, want %q, got %q", want, str)
	}
}
