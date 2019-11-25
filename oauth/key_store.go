package oauth

import (
	"bytes"
	"crypto/rsa"
	"crypto/sha256"
	"encoding/base64"
	"io/ioutil"

	"shanhu.io/base/httputil"
	"shanhu.io/misc/errcode"
	"shanhu.io/misc/rsautil"
)

func errUserNotFound(u string) error {
	return errcode.NotFoundf("user %q not found", u)
}

func parseKeys(bs []byte) ([]*PublicKey, error) {
	lines := bytes.Split(bs, []byte{'\n'})
	var keys []*PublicKey
	for _, line := range lines {
		line = bytes.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		k, err := newPublicKey(line)
		if err != nil {
			return nil, err
		}
		keys = append(keys, k)
	}
	return keys, nil
}

// PublicKey carries a public key.
type PublicKey struct {
	bytes []byte
	key   *rsa.PublicKey
	hash  []byte
}

func newPublicKey(bs []byte) (*PublicKey, error) {
	k, err := rsautil.ParsePublicKey(bs)
	if err != nil {
		return nil, err
	}

	h := sha256.Sum256(bs)

	return &PublicKey{
		bytes: bs,
		key:   k,
		hash:  h[:],
	}, nil
}

// Key returns the public key parsed from the bytes.
func (k *PublicKey) Key() *rsa.PublicKey { return k.key }

// HashStr returns the base64 encoding of the key hash.
func (k *PublicKey) HashStr() string {
	return base64.RawURLEncoding.EncodeToString(k.hash)
}

// KeyStore loads public keys for a user.
type KeyStore interface {
	Keys(user string) ([]*PublicKey, error)
}

// MemKeyStore is a storage of public keys in memory.
type MemKeyStore struct {
	keys map[string][]*PublicKey
}

// NewMemKeyStore creates a new empty key store.
func NewMemKeyStore() *MemKeyStore {
	return &MemKeyStore{keys: make(map[string][]*PublicKey)}
}

// Set sets the key for the given user.
func (s *MemKeyStore) Set(user string, keys []*PublicKey) {
	s.keys[user] = keys
}

// Keys returns the public keys for the given user.
func (s *MemKeyStore) Keys(user string) ([]*PublicKey, error) {
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
func (s *FileKeyStore) Keys(user string) ([]*PublicKey, error) {
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
	return parseKeys(bs)
}

// WebKeyStore is a storage of public keys backed by a web site.
type WebKeyStore struct {
	base   string
	client *httputil.Client
}

// NewWebKeyStore creates a new key store backed by a web site
// at the given base URL.
func NewWebKeyStore(base string) *WebKeyStore {
	return &WebKeyStore{
		base:   base,
		client: httputil.NewClient(base),
	}
}

// Keys returns the public keys of the given user.
func (s *WebKeyStore) Keys(user string) ([]*PublicKey, error) {
	bs, err := s.client.GetBytes(user)
	if err != nil {
		return nil, err
	}
	return parseKeys(bs)
}
