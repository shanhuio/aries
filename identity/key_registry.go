package identity

import (
	"io/ioutil"
	"path/filepath"
	"sync"

	"shanhu.io/misc/errcode"
	"shanhu.io/misc/rsautil"
)

func errUserNotFound(u string) error {
	return errcode.NotFoundf("user %q not found", u)
}

// KeyRegistry loads public keys for a user.
type KeyRegistry interface {
	Keys(user string) ([]*rsautil.PublicKey, error)
}

// MemKeyRegistry is a storage of public keys in memory.
type MemKeyRegistry struct {
	mu   sync.RWMutex
	keys map[string][]*rsautil.PublicKey
}

// NewMemKeyRegistry creates a new empty key store.
func NewMemKeyRegistry() *MemKeyRegistry {
	return &MemKeyRegistry{
		keys: make(map[string][]*rsautil.PublicKey),
	}
}

// Set sets the key for the given user.
func (r *MemKeyRegistry) Set(user string, keys []*rsautil.PublicKey) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.keys[user] = keys
}

// Keys returns the public keys for the given user.
func (r *MemKeyRegistry) Keys(user string) ([]*rsautil.PublicKey, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	keys, found := r.keys[user]
	if !found {
		return nil, errUserNotFound(user)
	}
	return keys, nil
}

// NewDirKeyRegistry creates a new keystore with public keys saved in
// files under a directory.
func NewDirKeyRegistry(dir string) (*MemKeyRegistry, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	m := NewMemKeyRegistry()
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		name := f.Name()
		if !IsSimpleName(name) {
			continue
		}
		bs, err := ioutil.ReadFile(filepath.Join(dir, name))
		if err != nil {
			return nil, errcode.Annotatef(err, "read key %q", name)
		}
		keys, err := rsautil.ParsePublicKeys(bs)
		if err != nil {
			return nil, errcode.Annotatef(err, "parse key %q", name)
		}
		m.Set(name, keys)
	}

	return m, nil
}

// FileKeyRegistry is a storage of public keys backed by a file system.
type FileKeyRegistry struct {
	keys map[string]string
}

// NewFileKeyRegistry creates a new key store given a key file
// map for each users that has a key.
func NewFileKeyRegistry(keys map[string]string) *FileKeyRegistry {
	return &FileKeyRegistry{keys: keys}
}

// Keys returns the public keys for the given user.
func (s *FileKeyRegistry) Keys(user string) ([]*rsautil.PublicKey, error) {
	if s.keys == nil {
		return nil, errUserNotFound(user)
	}

	f, found := s.keys[user]
	if !found {
		return nil, errUserNotFound(user)
	}
	bs, err := ioutil.ReadFile(f)
	if err != nil {
		return nil, err
	}
	return rsautil.ParsePublicKeys(bs)
}

// IsSimpleName checks if the user name is a simple one that is safe to
// fetch a key.
func IsSimpleName(user string) bool {
	for _, r := range user {
		if r >= 'a' && r <= 'z' {
			continue
		}
		if r >= '0' && r <= '9' {
			continue
		}
		if r == '~' || r == '-' || r == '_' {
			continue
		}
		return false
	}
	return true
}
