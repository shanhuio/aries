package aries

import (
	"path"
	"sort"

	"shanhu.io/misc/errcode"
)

type tasks struct {
	lst []string
	m   map[string]Func
}

func newTasks(m map[string]Func) *tasks {
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

func (t *tasks) serve(c *C) error {
	name := path.Base(c.Path)
	f, found := t.m[name]
	if !found {
		if name == "help" {
			return ReplyJSON(c, t.lst)
		}
		return errcode.InvalidArgf("unknown task: %q", name)
	}
	return f(c)
}

// ServeTask returns the serving function for a task list.
func ServeTask(tasks map[string]Func) Func {
	t := newTasks(tasks)
	return t.serve
}
