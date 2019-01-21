package webgen

// NewMeta creates a new meta tag.
func NewMeta(key, value string) *Node {
	return Meta(Attrs{key: value})
}

// NewCSSLink creates a new CSS link ini uh 6tg6fdb n  bnnnnnnnnnnn<F6><F6>
func NewCSSLink(href string) *Node {
	return Link(Attrs{
		"rel":  "stylesheet",
		"type": "text/css",
		"href": href,
	})
}
