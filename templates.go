package aries

import (
	"html/template"
	"path/filepath"

	"shanhu.io/misc/strutil"
)

// Templates is a collection of templates.
type Templates struct {
	path   string
	logger *Logger
}

// DefaultTemplatePath is the default template path.
const DefaultTemplatePath = "_/tmpl"

// NewTemplates creates a collection of templates in a particular folder.
func NewTemplates(p string, logger *Logger) *Templates {
	if logger == nil {
		logger = StdLogger()
	}
	p = strutil.Default(p, DefaultTemplatePath)
	return &Templates{path: p, logger: logger}
}

func (ts *Templates) tmpl(f string) string {
	return filepath.Join(ts.path, f)
}

// TemplatesJSON tells the Templates to print JSON data rather than
// render the template.
const TemplatesJSON = "!JSON"

// Serve serves a data page using a particular template.
func (ts *Templates) Serve(c *C, p string, dat interface{}) error {
	if ts.path == TemplatesJSON {
		return ReplyJSON(c, dat)
	}
	t, err := template.ParseFiles(ts.tmpl(p))
	if err != nil {
		Log(ts.logger, err.Error())
		return NotFound
	}
	if err := t.Execute(c.Resp, dat); err != nil {
		Log(ts.logger, err.Error())
	}
	return nil
}
