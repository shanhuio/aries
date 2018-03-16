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
