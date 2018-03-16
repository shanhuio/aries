package oauth

import (
	"io/ioutil"
)

// KeyStore loads a public key for a user.
type KeyStore interface {
	Key(user string) ([]byte, error)
}

// MemKeyStore is a storage of public keys in memory.
type MemKeyStore struct {
	keys map[string][]byte
}

// NewMemKeyStore creates a new empty key store.
func NewMemKeyStore() *MemKeyStore {
	return &MemKeyStore{keys: make(map[string][]byte)}
}

// Set sets the key for the given user.
func (s *MemKeyStore) Set(user string, k []byte) {
	cp := make([]byte, len(k))
	copy(cp, k)
	s.keys[user] = cp
}

// Key reads the key for the given user.
func (s *MemKeyStore) Key(user string) ([]byte, error) {
	bs, found := s.keys[user]
	if !found {
		return nil, nil
	}
	ret := make([]byte, len(bs))
	copy(ret, bs)
	return ret, nil
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
