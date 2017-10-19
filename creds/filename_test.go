package creds

import (
	"testing"
)

func TestFilename(t *testing.T) {
	o := func(from, to string) {
		got := Filename(from)
		if got != to {
			t.Errorf("Filename(%q) mapped to %q, want %q", from, got, to)
		}
	}

	o("shanhu.io", "shanhu-io")
	o("smallrepo.com", "smallrepo-com")
	o("localhost:3356", "localhost-3356")
	o("localhost:3335", "localhost-3335")
}
