package oauth

import (
	"io/ioutil"
	"net/url"
	"path/filepath"

	"shanhu.io/misc/errcode"
	"shanhu.io/misc/httputil"
	"shanhu.io/misc/rsautil"
)

func errUserNotFound(u string) error {
	return errcode.NotFoundf("user %q not found", u)
}

// KeyStore loads public keys for a user.
type KeyStore interface {
	Keys(user string) ([]*rsautil.PublicKey, error)
}

// MemKeyStore is a storage of public keys in memory.
type MemKeyStore struct {
	keys map[string][]*rsautil.PublicKey
}

// NewMemKeyStore creates a new empty key store.
func NewMemKeyStore() *MemKeyStore {
	return &MemKeyStore{
		keys: make(map[string][]*rsautil.PublicKey),
	}
}

// Set sets the key for the given user.
func (s *MemKeyStore) Set(user string, keys []*rsautil.PublicKey) {
	s.keys[user] = keys
}

// Keys returns the public keys for the given user.
func (s *MemKeyStore) Keys(user string) ([]*rsautil.PublicKey, error) {
	keys, found := s.keys[user]
	if !found {
		return nil, errUserNotFound(user)
	}
	return keys, nil
}

// FileKeyStore is a storage of public keys backed by a file system.
type FileKeyStore struct {
	keys map[string]string
}

// NewFileKeyStore creates a new key store given a key file
// map for each users that has a key.
func NewFileKeyStore(keys map[string]string) *FileKeyStore {
	return &FileKeyStore{keys: keys}
}

// Keys returns the public keys for the given user.
func (s *FileKeyStore) Keys(user string) ([]*rsautil.PublicKey, error) {
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

func simpleName(user string) bool {
	for _, r := range user {
		if r >= 'a' && r <= 'z' {
			continue
		}
		if r >= '0' && r <= '9' {
			continue
		}
		return false
	}
	return true
}

// DirKeyStore is a storage of pulic keys with all the keys saved in a
// directory.
type DirKeyStore struct {
	dir string
}

// NewDirKeyStore creates a new keystore with public keys saved in
// files under a directory.
func NewDirKeyStore(dir string) *DirKeyStore {
	return &DirKeyStore{dir: dir}
}

// Keys returns the public keys of the given user.
func (s *DirKeyStore) Keys(user string) ([]*rsautil.PublicKey, error) {
	if !simpleName(user) {
		return nil, errcode.InvalidArgf("unsupported user name: %q", user)
	}
	bs, err := ioutil.ReadFile(filepath.Join(s.dir, user))
	if err != nil {
		return nil, err
	}
	return rsautil.ParsePublicKeys(bs)
}

// WebKeyStore is a storage of public keys backed by a web site.
type WebKeyStore struct {
	client *httputil.Client
}

// NewWebKeyStore creates a new key store backed by a web site
// at the given base URL.
func NewWebKeyStore(base string) (*WebKeyStore, error) {
	client, err := httputil.NewClient(base)
	if err != nil {
		return nil, err
	}
	return &WebKeyStore{client: client}, nil
}

// Keys returns the public keys of the given user.
func (s *WebKeyStore) Keys(user string) ([]*rsautil.PublicKey, error) {
	if !simpleName(user) {
		return nil, errcode.InvalidArgf("unsupported user name: %q", user)
	}
	bs, err := s.client.GetBytes(user)
	if err != nil {
		return nil, err
	}
	return rsautil.ParsePublicKeys(bs)
}

// OpenKeyStore connects to a keystore based on the given URL string.
func OpenKeyStore(urlStr string) (KeyStore, error) {
	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	switch u.Scheme {
	case "http", "https":
		ks, err := NewWebKeyStore(urlStr)
		if err != nil {
			return nil, err
		}
		return ks, nil
	case "file", "":
		return NewDirKeyStore(u.Path), nil
	}
	return nil, errcode.InvalidArgf("unsupported url scheme: %q", u.Scheme)
}
