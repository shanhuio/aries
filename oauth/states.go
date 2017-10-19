package oauth

import (
	"net/http"
)

func stateCode(req *http.Request) (state, code string) {
	values := req.URL.Query()
	state = values.Get("state")
	if state != "" {
		code = values.Get("code")
	}
	return state, code
}
