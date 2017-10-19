package sitter

import (
	"fmt"
	"net/http"
	"sync"
)

type oneLine struct {
	mu   sync.RWMutex
	line string
}

func newOneLine(line string) *oneLine {
	return &oneLine{line: line}
}

func (s *oneLine) setLine(line string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.line = line
}

func (s *oneLine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	fmt.Fprint(w, s.line)
}
