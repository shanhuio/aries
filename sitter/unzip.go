package sitter

import (
	"archive/zip"
	"io"

	"shanhu.io/misc/tempfile"
	"shanhu.io/misc/ziputil"
)

func unzip(r io.Reader, dir string) error {
	tmp, err := tempfile.NewFile("", "sitter")
	if err != nil {
		return err
	}

	defer tmp.CleanUp()

	n, err := io.Copy(tmp, r)
	if err != nil {
		return err
	}
	zr, err := zip.NewReader(tmp, n)
	if err != nil {
		return err
	}
	if err := ziputil.UnzipDir(dir, zr, false); err != nil {
		return err
	}

	return tmp.CleanUp()
}
