package gometa

import (
	"fmt"
	"strings"
)

func isRepoPkg(repo, pkg string) bool {
	if repo == pkg {
		return true
	}
	return strings.HasPrefix(pkg, repo+"/")
}

func splitPkgName(pkg string) ([]string, error) {
	if pkg == "" {
		return nil, fmt.Errorf("empty package name")
	}

	parts := strings.Split(pkg, "/")
	for _, part := range parts {
		if part == "" {
			return nil, fmt.Errorf("invalid package name: %q", pkg)
		}
	}
	return parts, nil
}
