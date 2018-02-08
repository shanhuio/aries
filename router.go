package aries

import (
	"fmt"

	"shanhu.io/misc/trie"
)

// Router is a path router. Similar to mux, but routing base on
// a filesystem-like syntax.
type Router struct {
	index Func
	miss  Func

	trie  *trie.Trie
	funcs map[string]Func
}

// NewRouter creates a new router for filesystem like path routing.
func NewRouter() *Router {
	return &Router{
		trie:  trie.New(),
		funcs: make(map[string]Func),
	}
}

// Index sets the handler function for handling the index page when hitting
// this router, that is when hitting the root of it. One can only hit this
// route node when the path is ended with a slash '/'.
func (r *Router) Index(f Func) { r.index = f }

// Default sets a default handler for handling routes that does
// not hit anything in the routing tree.
func (r *Router) Default(f Func) { r.miss = f }

// Add adds a routing node into the routing tree
func (r *Router) Add(p string, f Func) error {
	if f == nil {
		panic("function is nil")
	}

	route := newRoute(p)
	if route.p == "" {
		panic("trying to add empty route, use Index() instead")
	}
	if r.funcs[route.p] != nil {
		return fmt.Errorf("path %s already assigned", route.p)
	}

	r.funcs[route.p] = f
	ok := r.trie.Add(route.routes, route.p)
	if !ok {
		panic("adding to trie failed")
	}

	return nil
}

func (r *Router) notFound(c *C) error {
	if r.miss == nil {
		return NotFound
	}
	return r.miss(c)
}

// Serve serves the incoming context. It returns Miss if the path hits
// nothing and Default() is not set.
func (r *Router) Serve(c *C) error {
	rel := c.route.rel(c.routePos)
	if rel == "" {
		if !c.route.isDir || r.index == nil {
			return r.notFound(c)
		}
		return r.index(c)
	}

	route := c.route.relRoute(c.routePos)
	p := r.trie.Find(route)
	f, ok := r.funcs[p]
	if !ok {
		return r.notFound(c)
	}

	// TODO: this is not right...
	return f(c)
}
