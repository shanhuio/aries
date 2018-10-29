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
	nodes map[string]*routerNode
}

type routerNode struct {
	f     Func
	isDir bool
}

// NewRouter creates a new router for filesystem like path routing.
func NewRouter() *Router {
	return &Router{
		trie:  trie.New(),
		nodes: make(map[string]*routerNode),
	}
}

// Index sets the handler function for handling the index page when hitting
// this router, that is when hitting the root of it. One can only hit this
// route node when the path is ended with a slash '/'.
func (r *Router) Index(f Func) { r.index = f }

// Default sets a default handler for handling routes that does
// not hit anything in the routing tree.
func (r *Router) Default(f Func) { r.miss = f }

// File adds a routing file node into the routing tree.
func (r *Router) File(p string, f Func) error {
	return r.add(p, &routerNode{f: f})
}

// Dir adds a routing directory node into the routing tree.
func (r *Router) Dir(p string, f Func) error {
	return r.add(p, &routerNode{f: f, isDir: true})
}

func (r *Router) add(p string, n *routerNode) error {
	if n.f == nil {
		panic("function is nil")
	}

	route := newRoute(p)
	if route.p == "" {
		panic("trying to add empty route, use Index() instead")
	}
	if r.nodes[route.p] != nil {
		return fmt.Errorf("path %s already assigned", route.p)
	}

	r.nodes[route.p] = n
	ok := r.trie.Add(route.routes, route.p)
	if !ok {
		panic("adding to trie failed")
	}

	return nil
}

func (r *Router) notFound(c *C) error {
	if r.miss == nil {
		return Miss
	}
	return r.miss(c)
}

// Serve serves the incoming context. It returns Miss if the path hits
// nothing and Default() is not set.
func (r *Router) Serve(c *C) error {
	rel := c.Rel()
	if rel == "" {
		if !c.PathIsDir() || r.index == nil {
			return r.notFound(c)
		}
		return r.index(c)
	}

	route := c.RelRoute()
	hitRoute, p := r.trie.Find(route)
	if p == "" {
		return r.notFound(c)
	}
	n := r.nodes[p]
	if n == nil {
		panic(fmt.Errorf("route function not found for %q", p))
	}

	c.ShiftRoute(len(hitRoute))
	if n.isDir || (c.Rel() == "" && !c.PathIsDir()) {
		return n.f(c)
	}
	return r.notFound(c)
}
