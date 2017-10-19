package sitter

import (
	"errors"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

type space struct {
	Name string
	Bin  string
	Dir  string

	cmd *exec.Cmd
}

func newSpace(name, dir string) (*space, error) {
	if err := os.MkdirAll(dir, 0700); err != nil {
		return nil, err
	}
	return &space{Name: name, Dir: dir}, nil
}

func (s *space) cleanUp() error { return os.RemoveAll(s.Dir) }

func (s *space) putFile(name string, r io.Reader, perm os.FileMode) error {
	p := filepath.Join(s.Dir, name)
	const mode = os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	f, err := os.OpenFile(p, mode, perm)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := io.Copy(f, r); err != nil {
		return err
	}

	return f.Close()
}

func (s *space) startAtAddr(addr string) error {
	return s.start([]string{s.Bin, "-addr=" + addr})
}

func (s *space) start(args []string) error {
	if s.cmd != nil {
		return errors.New("already running")
	}
	log.Println(s.Dir, "$ ", s.Bin, args)

	s.cmd = &exec.Cmd{
		Path:   s.Bin,
		Args:   args,
		Dir:    s.Dir,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	return s.cmd.Start()
}

func (s *space) stop() error {
	if s.cmd == nil || s.cmd.Process == nil {
		return nil
	}
	return s.cmd.Process.Kill()
}
