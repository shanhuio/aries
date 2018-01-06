package aries

import (
	"html/template"
	"log"
	"path/filepath"

	"shanhu.io/misc/errcode"
)

// Templates is a collection of templates.
type Templates struct {
	path string
}

// NewTemplates creates a collection of templates in a particular folder.
func NewTemplates(p string) *Templates {
	return &Templates{path: p}
}

func (ts *Templates) tmpl(f string) string {
	return filepath.Join(ts.path, f)
}

// Serve serves a webapp session with a particular template.
func (ts *Templates) Serve(c *C, p string, dat interface{}) error {
	t, err := template.ParseFiles(ts.tmpl(p))
	if err != nil {
		log.Println(err)
		return errcode.NotFoundf("page not found")
	}
	if err := t.Execute(c.Resp, dat); err != nil {
		log.Println(err)
	}
	return nil
}
