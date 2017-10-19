package aries

import (
	"testing"
)

func TestTrie(t *testing.T) {
	root := newTrieRoot()
	add := func(input string, added bool) {
		if added != root.add(input) {
			t.Errorf("expected dulplicated add = %v, got %v", added, !added)
		}
	}
	find := func(input, pref string, match bool) {
		rp, rm := trieFind(root, input)
		if rp != pref {
			t.Errorf("expected pref = %q, got %q", pref, rp)
		}
		if rm != match {
			t.Errorf("expected match = %v, got %v", match, rm)
		}
	}
	add("", false)
	find("a", "", false)
	add("axy45678", true)
	add("abc", true)
	add("axy", true)
	find("abc", "abc", true)
	find("abcd", "abc", false)
	find("a", "", false)
	find("ax", "", false)
	add("abc", false)
	add("a", true)
	find("a", "a", true)
	find("axy12", "axy", false)
	add("axy120", true)
	find("axy12", "axy", false)
	find("", "", true)
	add("ax456", true)
	add("ax4567", true)
	find("ax45", "a", false)
	find("dtc", "", false)
	find("axy45678", "axy45678", true)
}
