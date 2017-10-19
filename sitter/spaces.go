package sitter

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"sync"

	"shanhu.io/misc/jsonfile"
	"shanhu.io/misc/tempfile"
)

type state struct {
	Running string
	Spaces  map[string]*space
}

func newState() *state {
	return &state{
		Spaces: make(map[string]*space),
	}
}

type spaces struct {
	mu   sync.Mutex
	file string

	*state
}

func newSpaces(f string) *spaces {
	return &spaces{
		file:  f,
		state: newState(),
	}
}

func loadSpaces(f string) *spaces {
	s := newSpaces(f)
	err := jsonfile.Read(f, s.state)
	if err != nil && !os.IsNotExist(err) {
		log.Println(err)
	}
	return s
}

func (ss *spaces) writeState(w io.Writer) error {
	ss.mu.Lock()
	defer ss.mu.Unlock()
	enc := json.NewEncoder(w)
	return enc.Encode(ss.state)
}

func (ss *spaces) save() error {
	fout, err := tempfile.NewFile("", "sitter")
	if err != nil {
		return err
	}
	defer fout.CleanUp()

	if err := ss.writeState(fout); err != nil {
		return err
	}
	if err := fout.Close(); err != nil {
		return err
	}
	fout.SkipCleanUp = true
	if err := os.Rename(fout.Name, ss.file); err != nil {
		fout.Remove()
		return err
	}
	return nil
}

func (ss *spaces) get(name string) (*space, error) {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	ret, ok := ss.Spaces[name]
	if !ok {
		return nil, fmt.Errorf("space %q not found", name)
	}
	return ret, nil
}

func (ss *spaces) setBin(name, bin string) error {
	ss.mu.Lock()
	defer ss.mu.Unlock()
	ret, ok := ss.Spaces[name]
	if !ok {
		return fmt.Errorf("space %q not found", name)
	}
	ret.Bin = bin
	return nil
}

func (ss *spaces) add(s *space) error {
	ss.mu.Lock()
	defer ss.mu.Unlock()
	name := s.Name
	if _, ok := ss.Spaces[name]; ok {
		return fmt.Errorf("space %q already exist", name)
	}
	ss.Spaces[name] = s
	return nil
}

func (ss *spaces) setRunning(name string) error {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	switch name {
	case "", spaceInternal1, spaceInternal2:
		ss.Running = name
		return nil
	}

	if _, ok := ss.Spaces[name]; !ok {
		return fmt.Errorf("space %q not exist", name)
	}
	ss.Running = name
	return nil
}

func (ss *spaces) remove(name string) error {
	ss.mu.Lock()
	defer ss.mu.Unlock()

	if _, ok := ss.Spaces[name]; !ok {
		return fmt.Errorf("space %q not found", name)
	}
	if name == ss.Running {
		return fmt.Errorf("space %q is currently running", name)
	}

	delete(ss.Spaces, name)

	return nil
}
