package aries

import (
	"testing"

	"net/http/httptest"

	"shanhu.io/misc/errcode"
	"shanhu.io/misc/httputil"
)

func TestJSONCallString(t *testing.T) {
	const msg = "hello"
	const reply = "hi"

	f := func(c *C, in string) (string, error) {
		if in != msg {
			return "", errcode.InvalidArgf("not the right message")
		}
		return reply, nil
	}

	s := httptest.NewServer(JSONCall(f))
	defer s.Close()

	c := httputil.NewClient(s.URL)
	var str string
	if err := c.JSONCall("/", msg, &str); err != nil {
		t.Fatal(err)
	}

	if str != reply {
		t.Errorf("want %q, got %q", reply, str)
	}
}

func TestJSONFetchString(t *testing.T) {
	const reply = "hello"
	f := func(c *C) (string, error) { return reply, nil }

	s := httptest.NewServer(JSONCall(f))
	defer s.Close()

	c := httputil.NewClient(s.URL)
	var str string
	if err := c.JSONCall("/", nil, &str); err != nil {
		t.Fatal(err)
	}

	if str != reply {
		t.Errorf("want %q, got %q", reply, str)
	}
}

func TestJSONSendString(t *testing.T) {
	const msg = "hello"
	f := func(c *C, in string) error {
		if in != msg {
			return errcode.InvalidArgf("not the right message")
		}
		return nil
	}

	s := httptest.NewServer(JSONCall(f))
	defer s.Close()

	c := httputil.NewClient(s.URL)
	if err := c.JSONCall("/", msg, nil); err != nil {
		t.Fatal(err)
	}
}

func TestJSONCallStruct(t *testing.T) {
	type data struct {
		Message string
	}
	const msg = "hello"
	const reply = "hi"

	f := func(c *C, in *data) (*data, error) {
		if in.Message != msg {
			return nil, errcode.InvalidArgf("not the right message")
		}
		return &data{Message: reply}, nil
	}

	s := httptest.NewServer(JSONCall(f))
	defer s.Close()

	c := httputil.NewClient(s.URL)
	d := new(data)
	if err := c.JSONCall("/", &data{Message: msg}, d); err != nil {
		t.Fatal(err)
	}

	if d.Message != reply {
		t.Errorf("want %q, got %q", reply, d.Message)
	}
}
