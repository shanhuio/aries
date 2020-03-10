package oauth

import (
	"testing"

	"io/ioutil"
	"net/http/httptest"

	"shanhu.io/aries"
	"shanhu.io/misc/errcode"
	"shanhu.io/misc/rsautil"
)

func TestMemKeyStore(t *testing.T) {
	keyBytes, err := ioutil.ReadFile("testdata/keys/yumuzi")
	if err != nil {
		t.Fatal(err)
	}

	keys, err := rsautil.ParsePublicKeys(keyBytes)
	if err != nil {
		t.Fatal(err)
	}

	s := NewMemKeyStore()
	s.Set("h8liu", keys)

	got, err := s.Keys("h8liu")
	if err != nil {
		t.Fatal(err)
	}
	if len(keys) != len(got) {
		t.Errorf("want %d keys, got %d", len(keys), len(got))
	}

	for i, k := range got {
		h := k.HashStr()
		want := keys[i].HashStr()
		if h != want {
			t.Errorf("for key %d, want hash %q, got %q", i, want, h)
		}
	}
}

func testFileKeyStore(t *testing.T, ks KeyStore) {
	t.Helper()

	for _, test := range []struct {
		user   string
		hashes []string
	}{
		{"h8liu", []string{"zFet8qN1eNMvCQQqZRLy9Yxe-smJa8jmu30rOvBMeXw"}},
		{"yumuzi", []string{
			"Rxf8wK9cdKA6Zhn6KtVjSF3WUPLfnjbRlHuduSiOMsg",
			"zUcyOLg7_GzRTo4MDpyTnIxh6gqgGemUq0si_NjRXc4",
		}},
		{"xuduoduo", nil},
	} {
		t.Logf("test key for: %s", test.user)

		got, err := ks.Keys(test.user)
		if err != nil {
			if len(test.hashes) == 0 && errcode.IsNotFound(err) {
				continue
			}
			t.Fatal(err)
		}
		if len(got) != len(test.hashes) {
			t.Errorf("want %d keys, got %d", len(test.hashes), len(got))
			continue
		}

		for i, want := range test.hashes {
			if gotHash := got[i].HashStr(); want != gotHash {
				t.Errorf("key %d, want hash %q, got %q", i, want, gotHash)
			}
		}
	}
}

func TestFileKeyStore(t *testing.T) {
	s := NewFileKeyStore(map[string]string{
		"h8liu":  "testdata/keys/h8liu",
		"yumuzi": "testdata/keys/yumuzi",
	})

	testFileKeyStore(t, s)
}

func TestWebKeyStore(t *testing.T) {
	static := aries.NewStaticFiles("testdata/keys")
	s := httptest.NewServer(aries.Serve(static))
	defer s.Close()

	t.Log(s.URL)
	ks, err := NewWebKeyStore(s.URL)
	if err != nil {
		t.Fatal(err)
	}
	testFileKeyStore(t, ks)
}
