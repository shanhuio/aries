package gometa

import (
	"testing"
)

func TestSplitPkgName(t *testing.T) {
	for _, pkg := range []string{
		"",
		"/x",
		"x//y",
		"/",
		"a/b/c/",
	} {
		parts, err := splitPkgName(pkg)
		if err == nil {
			t.Errorf(
				"split package %q got parts %v, want error",
				pkg, parts,
			)
		}
	}
}
