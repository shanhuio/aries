// Package tasks provides a way to implement a set of simple subcommand-like
// API for user to invoke. Although a task is a POST request, the request often
// has no body content.
package tasks

import (
	"net/http"
	"sort"

	"shanhu.io/aries"
	"shanhu.io/misc/errcode"
)

type tasks struct {
	lst []string
	m   map[string]aries.Service
}

func newTasks(m map[string]aries.Service) *tasks {
	var lst []string
	for name := range m {
		lst = append(lst, name)
	}
	sort.Strings(lst)
	return &tasks{
		lst: lst,
		m:   m,
	}
}

func (t *tasks) serve(c *aries.C) error {
	name := c.Rel()

	if c.Req.Method != http.MethodPost {
		return errcode.NotFoundf("task %q must use POST", name)
	}

	f, found := t.m[name]
	if !found {
		if name == "help" {
			return aries.ReplyJSON(c, t.lst)
		}
		return errcode.InvalidArgf("unknown task: %q", name)
	}
	return f.Serve(c)
}

// Serve returns the serving function for a task list.
func Serve(tasks map[string]aries.Service) aries.Func {
	return newTasks(tasks).serve
}
