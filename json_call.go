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

var (
	errType     = reflect.TypeOf((*error)(nil)).Elem()
	contextType = reflect.TypeOf((*C)(nil))
)

func newJSONCall(f interface{}) (*jsonCall, error) {
	t := reflect.TypeOf(f)
	if k := t.Kind(); k != reflect.Func {
		return nil, fmt.Errorf("input is %s, not a function", k)
	}

	c := &jsonCall{f: reflect.ValueOf(f)}

	numIn := t.NumIn()
	if numIn == 0 || t.In(0) != contextType {
		return nil, fmt.Errorf("must use *aries.C as first arg")
	}

	if numIn == 1 {
		c.noRequest = true
	} else if numIn == 2 {
		c.req = t.In(1)
	} else {
		return nil, fmt.Errorf("invalid number of input: %d", numIn)
	}

	numOut := t.NumOut()
	if numOut == 1 {
		c.noResponse = true
		if got := t.Out(0); got != errType {
			return nil, fmt.Errorf("must return error, got %s", got)
		}
	} else if numOut == 2 {
		if got := t.Out(1); got != errType {
			return nil, fmt.Errorf("must return an error, got %s", got)
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
		if j.req.Kind() != reflect.Ptr {
			req := reflect.New(j.req)
			if err := UnmarshalJSONBody(c, req.Interface()); err != nil {
				return err
			}
			ret = j.f.Call([]reflect.Value{reflect.ValueOf(c), req.Elem()})
		} else {
			req := reflect.New(j.req.Elem())
			if err := UnmarshalJSONBody(c, req.Interface()); err != nil {
				return err
			}
			ret = j.f.Call([]reflect.Value{reflect.ValueOf(c), req})
		}
	} else {
		ret = j.f.Call([]reflect.Value{reflect.ValueOf(c)})
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
