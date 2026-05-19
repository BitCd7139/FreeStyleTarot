package service

import (
	"FreeStyleTarot/storage"
	"fmt"
	"regexp"
	"strings"
	"sync"

	"go.uber.org/zap"
)

type TarotRAGChunk struct {
	Key         string
	Title       string
	Description string
	Upright     string
	Reversed    string
}

var (
	tarotRAGIndex map[string]TarotRAGChunk
	ragOnce       sync.Once

	sectionHeaderRe = regexp.MustCompile(`(?m)^## (.+?) (Description|Upright|Reversed)\s*$`)
)

func loadTarotRAG() {
	ragOnce.Do(func() {
		tarotRAGIndex = make(map[string]TarotRAGChunk)

		raw, err := storage.Assets.ReadFile("tarot_rag_data.md")
		if err != nil {
			raw, err = storage.Assets.ReadFile("tarot_rag_data_example.md")
			if err != nil {
				zap.S().Errorw("Failed to read tarot RAG data", "error", err)
				return
			}
			zap.S().Warn("tarot_rag_data.md not found, using tarot_rag_data_example.md")
		}

		for _, part := range splitRAGChunks(string(raw)) {
			chunk, err := parseRAGChunk(part)
			if err != nil {
				zap.S().Warnw("Skipping invalid RAG chunk", "error", err)
				continue
			}
			tarotRAGIndex[chunk.Key] = chunk
		}
		zap.S().Infow("Loaded tarot RAG index", "cards", len(tarotRAGIndex))
	})
}

func splitRAGChunks(content string) []string {
	content = strings.TrimSpace(content)
	if content == "" {
		return nil
	}
	parts := strings.Split(content, "\n---")
	var chunks []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			chunks = append(chunks, p)
		}
	}
	return chunks
}

func parseRAGChunk(raw string) (TarotRAGChunk, error) {
	loc := sectionHeaderRe.FindAllStringSubmatchIndex(raw, -1)
	if len(loc) == 0 {
		return TarotRAGChunk{}, fmt.Errorf("no sections found")
	}

	var chunk TarotRAGChunk
	for i, match := range loc {
		title := strings.TrimSpace(raw[match[2]:match[3]])
		section := raw[match[4]:match[5]]
		start := match[1]
		end := len(raw)
		if i+1 < len(loc) {
			end = loc[i+1][0]
		}
		body := strings.TrimSpace(raw[start:end])

		switch section {
		case "Description":
			chunk.Title = title
			chunk.Description = body
		case "Upright":
			chunk.Upright = body
		case "Reversed":
			chunk.Reversed = body
		}
	}
	if chunk.Title == "" {
		return TarotRAGChunk{}, fmt.Errorf("missing card title")
	}
	chunk.Key = normalizeCardKey(chunk.Title)
	return chunk, nil
}

func normalizeCardKey(title string) string {
	var b strings.Builder
	for _, r := range strings.ToLower(title) {
		if r >= 'a' && r <= 'z' || r >= '0' && r <= '9' {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func resolveCardKey(name string) string {
	key := normalizeCardKey(name)
	if _, ok := tarotRAGIndex[key]; ok {
		return key
	}
	if alt, ok := cardKeyAliases[key]; ok {
		if _, ok := tarotRAGIndex[alt]; ok {
			return alt
		}
	}
	for alias, canonical := range cardKeyAliases {
		if canonical == key {
			if _, ok := tarotRAGIndex[alias]; ok {
				return alias
			}
		}
	}
	return key
}

// cardKeyAliases bridges slug variants between frontend and RAG headings.
var cardKeyAliases = map[string]string{
	"judgment": "judgement",
}

func lookupRAGChunk(cardName string) (TarotRAGChunk, bool) {
	loadTarotRAG()
	key := resolveCardKey(cardName)
	chunk, ok := tarotRAGIndex[key]
	return chunk, ok
}
