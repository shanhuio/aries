package aries

import (
	"net/http"

	"smallrepo.com/base/httputil"
)

// AltError is a short hand to httputil.AltError
func AltError(w http.ResponseWriter, err error, msg string, code int) {
	httputil.AltError(w, err, msg, code)
}
