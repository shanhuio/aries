package gometa

import (
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"shanhu.io/misc/pathutil"
)

func get(c *http.Client, pkg string) (*Repo, error) {
	url, err := url.Parse("https://" + pkg)
	if err != nil {
		return nil, err
	}
	url.RawQuery = "go-get=1"

	resp, err := c.Get(url.String())
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	repo, err := ParseGoImport(resp.Body, pkg)
	if err != nil {
		return nil, err
	}
	if err := resp.Body.Close(); err != nil {
		return nil, err
	}

	return repo, nil
}

func getRepo(c *http.Client, pkg string) (*Repo, error) {
	ret, err := get(c, pkg)
	if err != nil {
		return nil, err
	}

	check, err := get(c, ret.ImportRoot)
	if err != nil {
		return nil, err
	}
	if check.ImportRoot != ret.ImportRoot {
		return nil, fmt.Errorf(
			"repo path mismatch for %q, sub has %q, parent %q",
			pkg, ret.ImportRoot, check.ImportRoot,
		)
	}
	if check.VCSRoot != ret.VCSRoot {
		return nil, fmt.Errorf(
			"vcs mismatch for %q, sub has %q, parent %q",
			pkg, ret.VCSRoot, check.VCSRoot,
		)
	}

	return check, nil
}

// GetRepo gets the repo meta data for a particular package.
func GetRepo(c *http.Client, pkg string) (*Repo, error) {
	parts, err := pathutil.Split(pkg)
	if err != nil {
		return nil, err
	}

	domain := parts[0]
	switch domain {
	case "github.com", "bitbucket.org":
		if len(parts) < 3 {
			return nil, fmt.Errorf("cannot find repo for pkg: %q", pkg)
		}

		repoPath := path.Join(parts[:3]...)

		return &Repo{
			ImportRoot: repoPath,
			VCS:        "git",
			VCSRoot:    "https://" + repoPath,
		}, nil
	}

	for i, part := range parts {
		if strings.HasSuffix(part, ".git") {
			repoPath := path.Join(parts[:i+1]...)
			return &Repo{
				ImportRoot: repoPath,
				VCS:        "git",
				VCSRoot:    "https://" + repoPath,
			}, nil
		}
	}

	return get(c, pkg)
}
