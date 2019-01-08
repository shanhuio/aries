package aries

import (
	"fmt"
	"reflect"
)

type jsonCall struct {
	f          reflect.Value
	noRequest  bool
	noResponse bool
	req        reflect.Type
	resp       reflect.Type
}

var errType = reflect.TypeOf(error(nil))

func newJSONCall(f interface{}) (*jsonCall, error) {
	t := reflect.TypeOf(f)
	if k := t.Kind(); k != reflect.Func {
		return nil, fmt.Errorf("input is %s, not a function", k)
	}

	c := &jsonCall{f: reflect.ValueOf(f)}

	numIn := t.NumIn()
	if numIn == 0 {
		c.noRequest = true
	} else if numIn == 1 {
		c.req = t.In(0)
	} else {
		return nil, fmt.Errorf("invalid number of input: %d", numIn)
	}

	numOut := t.NumIn()
	if numOut == 1 {
		c.noResponse = true
		if t.Out(0) != errType {
			return nil, fmt.Errorf("must return an error")
		}
	} else if numOut == 2 {
		if t.Out(1) != errType {
			return nil, fmt.Errorf("must return an error")
		}
		c.resp = t.Out(0)
	} else {
		return nil, fmt.Errorf("invalid number of output: %d", numOut)
	}

	return c, nil
}

func (j *jsonCall) call(c *C) error {
	if m := c.Req.Method; m != "POST" {
		return fmt.Errorf("method is %q; must use POST", m)
	}

	var ret []reflect.Value
	if !j.noRequest {
		req := reflect.New(j.req)
		if err := UnmarshalJSONBody(c, req.Interface()); err != nil {
			return err
		}
		ret = j.f.Call([]reflect.Value{req})
	} else {
		ret = j.f.Call(nil)
	}

	var resp, err reflect.Value
	if !j.noResponse {
		resp = ret[0]
		err = ret[1]
	} else {
		err = ret[0]
	}

	if !err.IsNil() {
		return err.Interface().(error)
	}

	if j.noResponse {
		return nil
	}

	return ReplyJSON(c, resp.Interface())
}

// JSONCall wraps a function of form `func(req *RequestType) (resp
// *ResponseType, error)` into a JSON marshalled RPC call.
func JSONCall(f interface{}) Func {
	call, err := newJSONCall(f)
	if err != nil {
		panic(err)
	}
	return call.call
}
