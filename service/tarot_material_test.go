package service

import "testing"

func TestTruncateRunes(t *testing.T) {
	got := truncateRunes("abcdef", 3)
	if got != "abc…" {
		t.Fatalf("got %q", got)
	}
	if truncateRunes("ab", 10) != "ab" {
		t.Fatal("short string changed")
	}
}
