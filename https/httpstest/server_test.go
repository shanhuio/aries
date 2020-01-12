package httpstest

import (
	"testing"

	"io/ioutil"
	"net/http"

	"shanhu.io/aries"
)

func checkBody(t *testing.T, resp *http.Response, msg string) {
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read body: %s", err)
	}
	got := string(bs)
	if got != msg {
		t.Errorf("response body want: %q, got %q", msg, got)
	}
}

func checkGet(t *testing.T, c *http.Client, url, msg string) {
	resp, err := c.Get(url)
	if err != nil {
		t.Fatalf("get %s: %s", url, msg)
	}
	defer resp.Body.Close()

	checkBody(t, resp, msg)
}

func TestServer(t *testing.T) {
	const msg = "hello"
	s, err := NewServer(
		[]string{"test.shanhu.io"},
		aries.StringFunc(msg),
	)
	if err != nil {
		t.Fatalf("create server: %s", err)
	}

	c := s.Client()
	resp, err := c.Get("https://test.shanhu.io")
	if err != nil {
		t.Fatalf("get: %s", err)
	}
	defer resp.Body.Close()

	checkBody(t, resp, msg)
}

func TestDualServer(t *testing.T) {
	const msg = "hello"

	s, err := NewDualServer(
		[]string{"test.shanhu.io"},
		aries.StringFunc(msg),
	)
	if err != nil {
		t.Fatalf("create server: %s", err)
	}

	c := s.Client()
	checkGet(t, c, "https://test.shanhu.io", msg)
	checkGet(t, c, "http://test.shanhu.io", msg)
}
