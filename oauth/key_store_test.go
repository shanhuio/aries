package oauth

import (
	"testing"

	"bytes"
)

func TestMemKeyStore(t *testing.T) {
	k := []byte("my key")
	s := NewMemKeyStore()
	s.Set("h8liu", k)

	got, err := s.Key("h8liu")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(k, got) {
		t.Errorf("want %q, got %q", string(k), string(got))
	}

	got[0] = 'x'
	got, err = s.Key("h8liu")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(k, got) {
		t.Errorf("2nd time, want %q, got %q", string(k), string(got))
	}
}

func TestFileKeyStore(t *testing.T) {
	s := NewFileKeyStore(map[string]string{
		"h8liu": "testdata/h8liu.key",
		"yumuzi": "testdata/yumuzi.key",
	})

	for _, test := range []struct {
		user, key string
	} {
		{"h8liu", "h8\n"},
		{"yumuzi", "work?\n"},
		{"xuduoduo", ""},
	} {
		got, err := s.Key(test.user)
		if err != nil {
			t.Fatal(err)
		}
		if string(got) != test.key {
			t.Errorf("want %q, got %q", test.key, string(got))
		}
	}
}
