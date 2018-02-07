package aries

// Router is a path router. Similar to mux, but routing base on
// a filesystem-like syntax.
type Router struct {
	index Func
	miss  Func
}

// Index sets the handler function for handling the index page when hitting
// this router, that is when hitting the root of it. One can only hit this
// route node when the path is ended with a slash '/'.
func (r *Router) Index(f Func) { r.index = f }

// Default sets a default handler for handling routes that does
// not hit anything in the routing tree.
func (r *Router) Default(f Func) { r.miss = f }

// Node adds a routing node into the routing tree
func (r *Router) Node(p string, f Func) error {
	// TODO:
	return nil
}

// Serve serves the incoming context. It returns Miss if the path hits
// nothing and Default() is not set.
func (r *Router) Serve(c *C) error {
	// TODO:
	return Miss
}
