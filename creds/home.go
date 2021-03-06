package creds

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"

	"shanhu.io/misc/jsonutil"
)

const homeDir = ".shanhu"

// Home returns the directory for saving the credentials and config files.
func Home() (string, error) {
	u, err := user.Current()
	if err != nil {
		return "", err
	}
	return filepath.Join(u.HomeDir, homeDir), nil
}

// MakeHome creates the home directory if it does not exist.
func MakeHome() error {
	h, err := Home()
	if err != nil {
		return err
	}
	return os.MkdirAll(h, 700)
}

// HomeFile returns the path of a file under the home directory.
func HomeFile(f string) (string, error) {
	h, err := Home()
	if err != nil {
		return "", err
	}
	return filepath.Join(h, f), nil
}

// CheckPrivateFile checks if the file is of the right permission bits.
func CheckPrivateFile(f string) error {
	info, err := os.Stat(f)
	if err != nil {
		return err
	}
	mod := info.Mode() & 0777
	if mod != 0600 {
		return fmt.Errorf("file %q is not of perm 0600 but %#o", f, mod)
	}
	return err
}

// ReadPrivateFile reads the confent of a file. The file must be mode 0600.
func ReadPrivateFile(f string) ([]byte, error) {
	if err := CheckPrivateFile(f); err != nil {
		return nil, err
	}
	return ioutil.ReadFile(f)
}

// ReadHomeFile reads the content of a file under the home directory.
// The file must be mode 0600.
func ReadHomeFile(f string) ([]byte, error) {
	p, err := HomeFile(f)
	if err != nil {
		return nil, err
	}
	return ReadPrivateFile(p)
}

// WriteHomeFile updates a file under the home directory.
func WriteHomeFile(f string, bs []byte) error {
	p, err := HomeFile(f)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(p, bs, 0600)
}

// WriteHomeJSONFile updates a file under the home directory with a
// JSON marshallable blob.
func WriteHomeJSONFile(f string, v interface{}) error {
	buf := new(bytes.Buffer)
	if err := jsonutil.Fprint(buf, v); err != nil {
		return err
	}
	return WriteHomeFile(f, buf.Bytes())
}

// ReadHomeJSONFile reads a file under the home directory into a JSON
// marshallable structure.
func ReadHomeJSONFile(f string, v interface{}) error {
	bs, err := ReadHomeFile(f)
	if err != nil {
		return err
	}
	return json.Unmarshal(bs, v)
}
