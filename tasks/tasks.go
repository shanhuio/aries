package tasks

import (
	"path"
	"sort"

	"shanhu.io/aries"
	"shanhu.io/misc/errcode"
)

type tasks struct {
	lst []string
	m   map[string]aries.Func
}

func newTasks(m map[string]aries.Func) *tasks {
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
	name := path.Base(c.Path)
	f, found := t.m[name]
	if !found {
		if name == "help" {
			return aries.ReplyJSON(c, t.lst)
		}
		return errcode.InvalidArgf("unknown task: %q", name)
	}
	return f(c)
}

// Serve returns the serving function for a task list.
func Serve(tasks map[string]aries.Func) aries.Func {
	t := newTasks(tasks)
	return t.serve
}
