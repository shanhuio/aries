package webgen

import (
	"golang.org/x/net/html"
)

// Node wraps around an html node.
type Node struct{ *html.Node }

// Add appends more stuff into the node.
func (n *Node) Add(children ...interface{}) {
	addChildren(n.Node, children...)
}

func text(s string) *html.Node {
	return &html.Node{
		Type: html.TextNode,
		Data: s,
	}
}

// Text creates a text node.
func Text(s string) *Node { return &Node{text(s)} }

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
