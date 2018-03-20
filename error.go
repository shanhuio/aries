package aries

import (
	"shanhu.io/misc/errcode"
)

const nothingHere = "nothing here"

// Miss is returned when a mux or router does not
// hit anything in its path lookup.
var Miss error = errcode.NotFoundf(nothingHere)

// NotFound is a true not found error.
var NotFound error = errcode.NotFoundf(nothingHere)

// NeedSignIn is returned when sign in is required for visiting a particular
// page.
var NeedSignIn error = errcode.Unauthorizedf("please sign in")
