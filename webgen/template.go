package webgen

import (
	"html/template"
)

// Template makes an HTML template with the given body.
func Template(p *Page, body *Node) (*template.Template, error) {
	s, err := RenderString(p, body)
	if err != nil {
		return nil, err
	}
	return template.New("index").Parse(s)
}

// TemplateBody makes an HTML template with the given elements as the body.
func TemplateBody(children ...interface{}) (*template.Template, error) {
	return Template(nil, Body(children...))
}
