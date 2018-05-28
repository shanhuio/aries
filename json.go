package aries

import (
	"encoding/json"
	"log"

	"shanhu.io/misc/errcode"
)

// ReplyJSON replies a JSON marshable object over the response.
func ReplyJSON(c *C, v interface{}) error {
	bs, err := json.Marshal(v)
	if err != nil {
		return errcode.Internalf("response encode error: %s", err)
	}

	if _, err := c.Resp.Write(bs); err != nil {
		log.Println(err)
	}
	return nil
}

// UnmarshalJSONBody unmarshals the body into a JSON object.
func UnmarshalJSONBody(c *C, v interface{}) error {
	dec := json.NewDecoder(c.Req.Body)
	if err := dec.Decode(v); err != nil {
		return errcode.Add(errcode.InvalidArg, err)
	}
	return nil
}

// PrintJSON replies a JSON marshaable object over the reponse with
// pretty printing.
func PrintJSON(c *C, v interface{}) error {
	bs, err := json.MarshalIndent(v, "", "    ")
	if err != nil {
		return errcode.Internalf("response encode error: %s", err)
	}

	if _, err := c.Resp.Write(bs); err != nil {
		log.Println(err)
	}
	return nil
}
