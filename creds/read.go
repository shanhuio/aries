package creds

import (
	"log"
	"strings"
)

// ReadAPIKey reads in the secret file for shanhu deployments and API calls.
func ReadAPIKey() (string, error) {
	bs, err := ReadHomeFile("key")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(bs)), nil
}

// UseOrReadAPIKey returns s when s is not empty, otherwise it reads for the
// key. If read has an error, it returns an empty string and logs the error.
func UseOrReadAPIKey(s string) string {
	if s != "" {
		return s
	}

	ret, err := ReadAPIKey()
	if err != nil {
		log.Println(err)
		return ""
	}
	return ret
}
