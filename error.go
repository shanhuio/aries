package aries

import (
	"errors"

	"shanhu.io/misc/errcode"
)

// AltError logs the error and returns an alternative error with code.
func AltError(err error, code, s string) error {
	if err == nil {
		return nil
	}
	Log.Println(s, err)
	return errcode.Add(code, errors.New(s))
}

// AltInternal logs the error and returns an alternative internal error.
func AltInternal(err error, s string) error {
	return AltError(err, errcode.Internal, s)
}

// AltInvalidArg logs the error and returns an alternative invalid arg error.
func AltInvalidArg(err error, s string) error {
	return AltError(err, errcode.InvalidArg, s)
}

const nothingHere = "nothing here"

// Miss is returned when a mux or router does not
// hit anything in its path lookup.
var Miss error = errcode.NotFoundf(nothingHere)

// NotFound is a true not found error.
var NotFound error = errcode.NotFoundf(nothingHere)

// NeedSignIn is returned when sign in is required for visiting a particular
// page.
var NeedSignIn error = errcode.Unauthorizedf("please sign in")
