package aries

import (
	"net/http"

	"shanhu.io/misc/httputil"
)

// AltError is a short hand to httputil.AltError
func AltError(w http.ResponseWriter, err error, msg string, code int) {
	httputil.AltError(w, err, msg, code)
}
