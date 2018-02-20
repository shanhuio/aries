package aries

import (
	"testing"
)

func TestRoutePath(t *testing.T) {
	for _, test := range []struct {
		path string
		cleaned string
		size int
		isDir bool
	} {
		{"/", "", 0, true},
		{"/something", "/something", 1, false},
		{"/something/", "/something", 1, true},
		{"/a/b/c", "/a/b/c", 3, false},
		{"/a//c", "/a/c", 2, false},
		{"/////", "", 0, true},
	} {
		r := newRoute(test.path)
		got := r.path()
		if got != test.cleaned {
			t.Errorf(
				"clean route for %q, want %q, got %q",
				test.path, test.cleaned, got,
			)
		}

		size := r.size()
		if size != test.size {
			t.Errorf(
				"route size for %q, want %d, got %d",
				test.path, test.size, size,
			)
		}

		if r.isDir != test.isDir {
			t.Errorf(
				"route for %q, want isDir=%t, got %t",
				test.path, test.isDir, r.isDir,
			)
		}
	}
}
