package aries

import (
	"bytes"
	"strings"
)

type routePart struct {
	start, end int
}

type route struct {
	p     string
	parts []*routePart
	isDir bool
}

func newRoute(p string) *route {
	if p == "" {
		return new(route)
	}
	w := new(bytes.Buffer)
	n := len(p)
	isDir := p[n-1] == '/'

	splits := strings.Split(p, "/")
	var parts []*routePart
	for _, s := range splits {
		if len(s) == 0 {
			continue
		}
		w.WriteString("/")
		start := w.Len()
		w.WriteString(s)
		end := w.Len()
		parts = append(parts, &routePart{
			start: start,
			end:   end,
		})
	}

	return &route{
		p:     w.String(),
		parts: parts,
		isDir: isDir,
	}
}

func (r *route) dir(i int) string {
	return r.p[:r.parts[i].start]
}

func (r *route) current(i int) string {
	part := r.parts[i]
	return r.p[part.start:part.end]
}

func (r *route) rel(i int) string {
	return r.p[r.parts[i].start:]
}
