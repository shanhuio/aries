package webgen

import (
	"html/template"
)

// Template makes an html template with the given body.
func Template(name string, p *Page, body *Node) (*template.Template, error) {
	s, err := RenderString(p, body)
	if err != nil {
		return nil, err
	}
	return template.New(name).Parse(s)
}
