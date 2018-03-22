package gometa

import (
	"testing"

	"shanhu.io/aries"
)

type testServer struct {
	m *aries.Mux
}

func newTestServer() *testServer {
	return &testServer{
		m: NewGitMux("shanhu.io", map[string]string{
			"shanhu.io/repoa": "testdata/repoa",
			"shanhu.io/repob": "testdata/repob",
		}),
	}
}

func (s *testServer) Serve(c *aries.C) error {
	if IsGoGetRequest(c.Req) {
		return s.m.Serve(c)
	}
	return aries.NotFound
}

func TestGetRepo(t *testing.T) {

}
