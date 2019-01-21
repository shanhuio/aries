package webgen

import (
	"fmt"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// Element create a new element.
func Element(n interface{}, children ...interface{}) *Node {
	var ret *html.Node
	switch n := n.(type) {
	case atom.Atom:
		ret = &html.Node{
			Type:     html.ElementNode,
			DataAtom: n,
			Data:     n.String(),
		}
	case string:
		ret = &html.Node{
			Type: html.ElementNode,
			Data: n,
		}
	default:
		panic(fmt.Sprintf("unknown element %T", n))
	}

	addChildren(ret, children...)
	return &Node{ret}
}

func bind(a atom.Atom) func(c ...interface{}) *Node {
	return func(c ...interface{}) *Node { return Element(a, c...) }
}

// Shorthand element creators.
var (
	Body       = bind(atom.Body)
	Div        = bind(atom.Div)
	Span       = bind(atom.Span)
	Title      = bind(atom.Title)
	Meta       = bind(atom.Meta)
	Head       = bind(atom.Head)
	HTML       = bind(atom.Html)
	P          = bind(atom.P)
	Pre        = bind(atom.Pre)
	Blockquote = bind(atom.Blockquote)
	Strong     = bind(atom.Strong)
	Em         = bind(atom.Em)
	H1         = bind(atom.H1)
	H2         = bind(atom.H2)
	H3         = bind(atom.H3)
	H4         = bind(atom.H3)
	H5         = bind(atom.H3)
	H6         = bind(atom.H3)
	A          = bind(atom.A)
)

// NewHTML creates a new blank HTML element with the specified language.
func NewHTML(lang string) *Node {
	if lang == "" {
		return HTML()
	}
	return HTML(Attrs{"lang": lang})
}

// NewHTMLEnglish creates a new English HTML element.
func NewHTMLEnglish() *Node { return NewHTML("en") }

// NewHTMLChinese creates a new Chinese HTML element.
func NewHTMLChinese() *Node { return NewHTML("zh") }

// NewMeta create a new meta tag.
func NewMeta(key, value string) *Node {
	return Meta(Attrs{key: value})
}

// NewLink creates a new web link.
func NewLink(href string, children ...interface{}) *Node {
	var stuff []interface{}
	stuff = append(stuff, Attrs{"href": href})
	stuff = append(stuff, children...)
	return A(stuff...)
}
