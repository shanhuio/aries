package aries

import (
	"fmt"

	"shanhu.io/misc/errcode"
	"shanhu.io/misc/trie"
)

// Router is a path router. Similar to mux, but routing base on
// a filesystem-like syntax.
type Router struct {
	index Service
	miss  Service

	trie  *trie.Trie
	nodes map[string]*routerNode
}

type routerNode struct {
	s      Service
	isDir  bool
	method string
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

// MethodFile adds a routing file node into the routing tree that accepts
// only the given method.
func (r *Router) MethodFile(m, p string, f Func) error {
	return r.add(p, &routerNode{s: f, method: m})
}

// File adds a routing file node into the routing tree.
func (r *Router) File(p string, f Func) error {
	return r.MethodFile("", p, f)
}

// Get adds a routing file node into the routing tree that handles GET
// requests.
func (r *Router) Get(p string, f Func) error {
	return r.MethodFile("GET", p, f)
}

// JSONCall adds a JSON marshalled POST based RPC call node into the routing
// tree. The function must be in the form of
// `func(c *aries.C, req *RequestType) (resp *ResponseType, error)`,
// where RequestType
// and ResponseType are both JSON marshallable.
func (r *Router) JSONCall(p string, f interface{}) error {
	return r.File(p, JSONCall(f))
}

// JSONCallMust is the same as JSONCall, but panics if there is an error.
func (r *Router) JSONCallMust(p string, f interface{}) {
	if err := r.JSONCall(p, f); err != nil {
		panic(err)
	}
}

// Call is an alias of JSONCallMust
func (r *Router) Call(p string, f interface{}) {
	r.JSONCallMust(p, f)
}

// Dir adds a routing directory node into the routing tree.
func (r *Router) Dir(p string, f Func) error {
	return r.DirService(p, f)
}

// DirService adds a service into the router tree under a directory node.
func (r *Router) DirService(p string, s Service) error {
	return r.add(p, &routerNode{s: s, isDir: true})
}

func (r *Router) add(p string, n *routerNode) error {
	if n.s == nil {
		panic("function is nil")
	}

	route := newRoute(p)
	if route.p == "" {
		panic("trying to add empty route, use Index() instead")
	}
	if r.nodes[route.p] != nil {
		return errcode.InvalidArgf("path %s already assigned", route.p)
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
	return r.miss.Serve(c)
}

// Serve serves the incoming context. It returns Miss if the path hits
// nothing and Default() is not set.
func (r *Router) Serve(c *C) error {
	rel := c.Rel()
	if rel == "" {
		if r.index == nil {
			return r.notFound(c)
		}
		return r.index.Serve(c)
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
		m := c.Req.Method
		if n.method != "" && m != n.method {
			return errcode.InvalidArgf("unsupported method: %q", m)
		}
		return n.s.Serve(c)
	}
	return r.notFound(c)
}
