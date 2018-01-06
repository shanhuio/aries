package aries

import (
	"errors"
	"log"

	"shanhu.io/misc/errcode"
)

// AltError logs the error and returns an alternative error with code.
func AltError(err error, code, s string) error {
	if err == nil {
		return nil
	}
	log.Println(s, err)
	return errcode.Add(code, errors.New(s))
}

// AltInternal logs the error and returns an alternative internal error.
func AltInternal(err error, s string) error {
	return AltError(err, errcode.Internal, s)
}

// AltInvalidArg logs the error and returns an alternative invalid arg error.
func (c *C) AltInvalidArg(err error, s string) error {
	return AltError(err, errcode.InvalidArg, s)
}
