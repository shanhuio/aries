package oauth

import (
	"io/ioutil"
)

// KeyStore loads a public key for a user.
type KeyStore interface {
	Key(user string) ([]byte, error)
}

// FileKeyStore is a storage of public keys.
type FileKeyStore struct {
	keys map[string]string
}

// NewFileKeyStore creates a new key store given a key file
// map for each users that has a key.
func NewFileKeyStore(keys map[string]string) *FileKeyStore {
	return &FileKeyStore{keys: keys}
}

// Key reads the key for the given user.
func (s *FileKeyStore) Key(user string) ([]byte, error) {
	if s.keys == nil {
		return nil, nil
	}

	f, found := s.keys[user]
	if !found {
		return nil, nil
	}
	return ioutil.ReadFile(f)
}
