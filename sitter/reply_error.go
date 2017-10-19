package sitter

import (
	"net/http"

	"shanhu.io/misc/httputil"
)

func replyError(w http.ResponseWriter, err error) bool {
	if err == nil {
		return false
	}

	httputil.AltError(w, err, err.Error(), 400)
	return true
}
