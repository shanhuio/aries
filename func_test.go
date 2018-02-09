package aries

import (
	"testing"

	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
)

func httpGetString(c *http.Client, url string) (string, error) {
	resp, err := c.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("status code: %d", resp.StatusCode)
	}
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(bs), nil
}

func httpGetCode(c *http.Client, url string) (int, error) {
	resp, err := c.Get(url)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	return resp.StatusCode, nil
}

func TestFunc(t *testing.T) {
	const msg = "hello"
	f := func(c *C) error {
		fmt.Fprint(c.Resp, msg)
		return nil
	}
	s := httptest.NewServer(Func(f))
	defer s.Close()

	got, err := httpGetString(s.Client(), s.URL)
	if err != nil {
		t.Error(err)
		return
	}
	if got != msg {
		t.Errorf("want %q in response, got %s", msg, got)
	}
}

func TestFuncHTTPS(t *testing.T) {
	const msg = "hello"
	f := func(c *C) error {
		fmt.Fprint(c.Resp, msg)
		return nil
	}
	s := httptest.NewTLSServer(Func(f))
	defer s.Close()

	got, err := httpGetString(s.Client(), s.URL)
	if err != nil {
		t.Error(err)
		return
	}
	if got != msg {
		t.Errorf("want %q in response, got %s", msg, got)
	}
}
