package webgen

// NewMeta create a new meta tag.
func NewMeta(key, value string) *Node {
	return Meta(Attrs{key: value})
}
