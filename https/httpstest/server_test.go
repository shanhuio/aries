package httpstest

import (
	"testing"

	"fmt"
	"io/ioutil"
	"net/http"
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

func TestServer(t *testing.T) {
	const msg = "hello"
	s, err := NewServer(
		[]string{"test.shanhu.io"},
		http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			fmt.Fprint(w, msg)
		}),
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
		http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			fmt.Fprint(w, msg)
		}),
	)
	if err != nil {
		t.Fatalf("create server: %s", err)
	}

	c := s.Client()
	resp1, err := c.Get("https://test.shanhu.io")
	if err != nil {
		t.Fatalf("get: %s", err)
	}
	defer resp1.Body.Close()

	checkBody(t, resp1, msg)

	resp2, err := c.Get("http://test.shanhu.io")
	if err != nil {
		t.Fatalf("get: %s", err)
	}
	defer resp2.Body.Close()

	checkBody(t, resp2, msg)
}
