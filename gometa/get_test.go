package gometa

import (
	"testing"

	"fmt"
	"reflect"
	"strings"

	"shanhu.io/aries"
	"shanhu.io/aries/ariestest"
)

type testServer struct {
	m *aries.Mux
}

func newTestServer() *testServer {
	return &testServer{
		m: NewGitMux("shanhu.io", map[string]string{
			"repoa": "repo/a",
			"repob": "repo/b",
		}),
	}
}

func serveFakeBitBucket(c *aries.C) error {
	if strings.Index(c.Path, "repod") >= 0 {
		fmt.Fprintln(c.Resp, `{"scm":"git"}`)
	} else if strings.Index(c.Path, "repo-hg") >= 0 {
		fmt.Fprintln(c.Resp, `{"scm":"hg"}`)
	}
	return aries.Miss
}

func (s *testServer) Serve(c *aries.C) error {
	if IsGoGetRequest(c.Req) {
		return s.m.Serve(c)
	}
	return aries.NotFound
}

func TestGetRepo(t *testing.T) {
	s, err := ariestest.HTTPSServers(map[string]aries.Service{
		"shanhu.io":         newTestServer(),
		"api.bitbucket.org": aries.Func(serveFakeBitBucket),
	})
	if err != nil {
		t.Fatal(err)
	}
	defer s.Close()

	repoa := NewGitRepo("shanhu.io/repoa", "repo/a")
	repob := NewGitRepo("shanhu.io/repob", "repo/b")
	repoc := NewGitRepo(
		"github.com/h8liu/repoc",
		"https://github.com/h8liu/repoc",
	)
	repod := NewGitRepo(
		"bitbucket.org/h8liu/repod",
		"https://bitbucket.org/h8liu/repod",
	)
	repoHG := &Repo{
		ImportRoot: "bitbucket.org/h8liu/repo-hg",
		VCS:        "hg",
		VCSRoot:    "https://bitbucket.org/h8liu/repo-hg",
	}

	c := s.Client()
	for _, test := range []struct {
		repo string
		want *Repo
	}{
		{"shanhu.io/repoa", repoa},
		{"shanhu.io/repob", repob},
		{"shanhu.io/repob/subpackage", repob},
		{"github.com/h8liu/repoc", repoc},
		{"github.com/h8liu/repoc/xxx", repoc},
		{"bitbucket.org/h8liu/repod/xxx", repod},
		{"bitbucket.org/h8liu/repo-hg/xxx", repoHG},
	} {
		repo, err := GetRepo(c, test.repo)
		if err != nil {
			t.Errorf("get repo %q, got error %s", test.repo, err)
		} else if !reflect.DeepEqual(repo, test.want) {
			t.Errorf(
				"get repo %q, got %v, want %v",
				test.repo, repo, test.want,
			)
		}
	}

	for _, url := range []string{
		"shanhu.io",
		"smlrepo.com/xxx",
	} {
		repo, err := GetRepo(c, url)
		if err == nil {
			t.Errorf("get repo %q, want error, got %v", url, repo)
		}
	}
}
