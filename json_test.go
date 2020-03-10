package aries

import (
	"testing"

	"net/http/httptest"

	"shanhu.io/misc/errcode"
	"shanhu.io/misc/httputil"
)

func TestJSONString(t *testing.T) {
	const msg = "hello"
	const reply = "hi"

	f := func(c *C) error {
		var s string
		if err := UnmarshalJSONBody(c, &s); err != nil {
			return err
		}
		if s != msg {
			return errcode.InvalidArgf("not the right message")
		}
		return ReplyJSON(c, reply)
	}

	s := httptest.NewServer(Func(f))
	defer s.Close()

	c := httputil.NewClientMust(s.URL)
	var str string
	if err := c.JSONCall("/", msg, &str); err != nil {
		t.Fatal(err)
	}

	if str != reply {
		t.Errorf("want %q, got %q", reply, str)
	}
}

func TestJSONStruct(t *testing.T) {
	type data struct {
		Message string
	}
	const msg = "hello"
	const reply = "hi"

	f := func(c *C) error {
		d := new(data)
		if err := UnmarshalJSONBody(c, d); err != nil {
			return err
		}
		if d.Message != msg {
			return errcode.InvalidArgf("not the right message")
		}
		return ReplyJSON(c, &data{Message: reply})
	}

	s := httptest.NewServer(Func(f))
	defer s.Close()

	c := httputil.NewClientMust(s.URL)
	d := new(data)
	if err := c.JSONCall("/", &data{Message: msg}, d); err != nil {
		t.Fatal(err)
	}

	if d.Message != reply {
		t.Errorf("want %q, got %q", reply, d.Message)
	}
}
