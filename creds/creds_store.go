package creds

import (
	"shanhu.io/misc/errcode"
)

type credsStore interface {
	read() (*Creds, error)
	write(c *Creds) error
}

type homeCredsStore struct {
	file string
}

func newHomeCredsStore(server string) *homeCredsStore {
	f := Filename(server) + ".json"
	return &homeCredsStore{file: f}
}

func (s *homeCredsStore) read() (*Creds, error) {
	creds := &Creds{}
	if err := ReadHomeJSONFile(s.file, creds); err != nil {
		return nil, errcode.FromOS(err)
	}
	return creds, nil
}

func (s *homeCredsStore) write(c *Creds) error {
	return errcode.FromOS(WriteHomeJSONFile(s.file, c))
}
