package creds

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
		return nil, err
	}
	return creds, nil
}

func (s *homeCredsStore) write(c *Creds) error {
	return WriteHomeJSONFile(s.file, c)
}
