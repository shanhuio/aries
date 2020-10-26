package identity

import (
	"io/ioutil"
	"path/filepath"

	"shanhu.io/misc/errcode"
	"shanhu.io/misc/rsautil"
)

// DirKeyRegistry is a storage of pulic keys with all the keys saved in a
// directory.
type DirKeyRegistry struct {
	dir string
}

// NewDirKeyRegistry creates a new keystore with public keys saved in
// files under a directory.
func NewDirKeyRegistry(dir string) (*DirKeyRegistry, error) {
	return &DirKeyRegistry{dir: dir}, nil
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
