package identity

import (
	"io/ioutil"
	"net/url"
	"path/filepath"
	"sync"

	"shanhu.io/misc/errcode"
	"shanhu.io/misc/httputil"
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
		if r == '~' {
			continue
		}
		return false
	}
	return true
}

// DirKeyRegistry is a storage of pulic keys with all the keys saved in a
// directory.
type DirKeyRegistry struct {
	dir string
}

// NewDirKeyRegistry creates a new keystore with public keys saved in
// files under a directory.
func NewDirKeyRegistry(dir string) *DirKeyRegistry {
	return &DirKeyRegistry{dir: dir}
}

// Keys returns the public keys of the given user.
func (s *DirKeyRegistry) Keys(user string) ([]*rsautil.PublicKey, error) {
	if !IsSimpleName(user) {
		return nil, errcode.InvalidArgf("unsupported user name: %q", user)
	}
	bs, err := ioutil.ReadFile(filepath.Join(s.dir, user))
	if err != nil {
		return nil, err
	}
	return rsautil.ParsePublicKeys(bs)
}

// WebKeyRegistry is a storage of public keys backed by a web site.
type WebKeyRegistry struct {
	client *httputil.Client
}

// NewWebKeyRegistry creates a new key store backed by a web site
// at the given base URL.
func NewWebKeyRegistry(base *url.URL) *WebKeyRegistry {
	client := &httputil.Client{Server: base}
	return &WebKeyRegistry{client: client}
}

// Keys returns the public keys of the given user.
func (s *WebKeyRegistry) Keys(user string) ([]*rsautil.PublicKey, error) {
	if !IsSimpleName(user) {
		return nil, errcode.InvalidArgf("unsupported user name: %q", user)
	}
	bs, err := s.client.GetBytes(user)
	if err != nil {
		return nil, err
	}
	return rsautil.ParsePublicKeys(bs)
}

// OpenKeyRegistry connects to a keystore based on the given URL string.
func OpenKeyRegistry(urlStr string) (KeyRegistry, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	switch u.Scheme {
	case "http", "https":
		u, err := url.Parse(urlStr)
		if err != nil {
			return nil, err
		}
		return NewWebKeyRegistry(u), nil
	case "file", "":
		return NewDirKeyRegistry(u.Path), nil
	}
	return nil, errcode.InvalidArgf("unsupported url scheme: %q", u.Scheme)
}
