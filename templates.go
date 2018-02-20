package aries

import (
	"html/template"
	"path/filepath"

	"shanhu.io/misc/strutil"
)

// Templates is a collection of templates.
type Templates struct {
	path string
}

// DefaultTemplatePath is the default template path.
const DefaultTemplatePath = "_/tmpl"

// NewTemplates creates a collection of templates in a particular folder.
func NewTemplates(p string) *Templates {
	p = strutil.Default(p, DefaultTemplatePath)
	return &Templates{path: p}
}

func (ts *Templates) tmpl(f string) string {
	return filepath.Join(ts.path, f)
}

// Serve serves a webapp session with a particular template.
func (ts *Templates) Serve(c *C, p string, dat interface{}) error {
	t, err := template.ParseFiles(ts.tmpl(p))
	if err != nil {
		c.Log.Println(err)
		return NotFound
	}
	if err := t.Execute(c.Resp, dat); err != nil {
		c.Log.Println(err)
	}
	return nil
}
