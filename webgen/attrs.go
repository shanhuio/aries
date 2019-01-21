package webgen

import (
	"sort"
	"strings"

	"golang.org/x/net/html"
)

// Attrs is an attribute map.
type Attrs map[string]string

func setAttrs(node *html.Node, attrs Attrs) {
	var keys []string
	for k := range attrs {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		node.Attr = append(node.Attr, html.Attribute{
			Key: k,
			Val: attrs[k],
		})
	}
	return
}

// Class is the class attribute for a div.
type Class []string

func setClass(node *html.Node, cls Class) {
	if len(cls) == 0 {
		return
	}

	node.Attr = append(node.Attr, html.Attribute{
		Key: "class",
		Val: strings.Join([]string(cls), " "),
	})
}
