package service

import (
	"strings"
	"testing"
)

func TestParseRAGChunk(t *testing.T) {
	raw := `## Two Of Swords Description

Desc body here.

## Two Of Swords Upright

Upright body.

## Two Of Swords Reversed

Reversed body.`

	chunk, err := parseRAGChunk(raw)
	if err != nil {
		t.Fatal(err)
	}
	if chunk.Key != "twoofswords" {
		t.Fatalf("key = %q, want twoofswords", chunk.Key)
	}
	if !strings.Contains(chunk.Description, "Desc body") {
		t.Fatalf("description: %q", chunk.Description)
	}
	if !strings.Contains(chunk.Upright, "Upright body") {
		t.Fatalf("upright: %q", chunk.Upright)
	}
	if !strings.Contains(chunk.Reversed, "Reversed body") {
		t.Fatalf("reversed: %q", chunk.Reversed)
	}
}

func TestNormalizeCardKey(t *testing.T) {
	if got := normalizeCardKey("The Empress"); got != "theempress" {
		t.Fatalf("got %q", got)
	}
}
