package webgen

import (
	"fmt"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func addChildren(n *html.Node, children ...interface{}) {
	for _, child := range children {
		switch c := child.(type) {
		case Attrs:
			setAttrs(n, c)
		case string:
			n.AppendChild(text(c))
		case *html.Node:
			n.AppendChild(c)
		case *Node:
			n.AppendChild(c.Node)
		}
	}
}

// Node wraps around an html node.
type Node struct{ *html.Node }

// Text creates a text node.
func Text(s string) *Node { return &Node{text(s)} }

func text(s string) *html.Node {
	return &html.Node{
		Type: html.TextNode,
		Data: s,
	}
}

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

// Add appends more stuff into the node.
func (n *Node) Add(children ...interface{}) {
	addChildren(n.Node, children...)
}

func bind(a atom.Atom) func(c ...interface{}) *Node {
	return func(c ...interface{}) *Node { return Element(a, c...) }
}

// Shorthand element creators.
var (
	Body   = bind(atom.Body)
	Div    = bind(atom.Div)
	Span   = bind(atom.Span)
	Title  = bind(atom.Title)
	Meta   = bind(atom.Meta)
	Head   = bind(atom.Head)
	P      = bind(atom.P)
	Strong = bind(atom.Strong)
	Em     = bind(atom.Em)
	H1     = bind(atom.H1)
	H2     = bind(atom.H2)
	H3     = bind(atom.H3)
	H4     = bind(atom.H3)
	H5     = bind(atom.H3)
	H6     = bind(atom.H3)
	A      = bind(atom.A)
)

// NewMeta create a new meta tag.
func NewMeta(key, value string) *Node {
	return Meta(Attrs{key: value})
}
