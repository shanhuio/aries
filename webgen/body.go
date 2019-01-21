package webgen

// NewLink creates a new web link.
func NewLink(href string, children ...interface{}) *Node {
	var stuff []interface{}
	stuff = append(stuff, Attrs{"href": href})
	stuff = append(stuff, children...)
	return A(stuff...)
}
