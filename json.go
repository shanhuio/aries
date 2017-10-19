package aries

import (
	"encoding/json"
	"log"
	"net/http"
)

// ReplyJSON replies a JSON marshable object over the response.
func ReplyJSON(c *C, v interface{}) error {
	bs, err := json.Marshal(v)
	if err != nil {
		log.Println(err)
		http.Error(c.Resp, "respond encode error", 400)
		return err
	}

	if _, err := c.Resp.Write(bs); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// UnmarshalJSONBody unmarshals the body into a JSON object.
func UnmarshalJSONBody(c *C, v interface{}) error {
	dec := json.NewDecoder(c.Req.Body)
	return dec.Decode(v)
}
