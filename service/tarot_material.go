package service

import (
	"FreeStyleTarot/model/request"
	"errors"
	"fmt"
	"sort"
	"strings"
)

// BuildSpreadMaterials assembles per-card markdown for the drawn spread.
func BuildSpreadMaterials(cards []request.CardInfo) (string, error) {
	loadKnowledgeBase()
	loadTarotRAG()

	if len(cards) == 0 {
		return "", errors.New("no cards in spread")
	}

	sorted := make([]request.CardInfo, len(cards))
	copy(sorted, cards)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Order < sorted[j].Order
	})

	var b strings.Builder
	b.WriteString("## Spread knowledge base\n\n")
	for _, card := range sorted {
		block, err := buildCardMaterial(card)
		if err != nil {
			return "", err
		}
		b.WriteString(block)
		b.WriteString("\n")
	}
	return strings.TrimSpace(b.String()), nil
}

// BuildSpreadBriefMaterials is a condensed knowledge digest for persona_advisor (Agent B).
func BuildSpreadBriefMaterials(cards []request.CardInfo) (string, error) {
	loadKnowledgeBase()
	loadTarotRAG()

	if len(cards) == 0 {
		return "", errors.New("no cards in spread")
	}

	maxRunes := briefMaterialMaxRunes
	if maxRunes <= 0 {
		maxRunes = 500
	}

	sorted := make([]request.CardInfo, len(cards))
	copy(sorted, cards)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Order < sorted[j].Order
	})

	var b strings.Builder
	b.WriteString("## Spread reference (digest)\n\n")
	for _, card := range sorted {
		block, err := buildCardBriefMaterial(card, maxRunes)
		if err != nil {
			return "", err
		}
		b.WriteString(block)
		b.WriteString("\n")
	}
	return strings.TrimSpace(b.String()), nil
}

func buildCardBriefMaterial(card request.CardInfo, maxRunes int) (string, error) {
	detail, hasKB := tarotMap[card.Name]
	if !hasKB {
		if alt, ok := cardKeyAliases[card.Name]; ok {
			detail, hasKB = tarotMap[alt]
		}
	}
	if !hasKB {
		return "", fmt.Errorf("card %q not found in knowledge base", card.Name)
	}

	reversed := strings.EqualFold(card.Orientation, "reversed")
	state := "正位"
	keyword := detail.Upright
	if reversed {
		state = "逆位"
		keyword = detail.Reversed
	}

	rag, hasRAG := lookupRAGChunk(card.Name)
	displayName := detail.Name
	if rag.Title != "" {
		displayName = rag.Title
	}

	desc := detail.Description
	interpret := ""
	if hasRAG {
		if rag.Description != "" {
			desc = rag.Description
		}
		interpret = rag.Upright
		if reversed {
			interpret = rag.Reversed
		}
	}

	var b strings.Builder
	fmt.Fprintf(&b, "### %s · %s", displayName, state)
	if card.Meaning != "" {
		fmt.Fprintf(&b, " · 阵位：%s", card.Meaning)
	}
	b.WriteString("\n\n")
	fmt.Fprintf(&b, "**关键词**: %s\n\n", keyword)
	if desc != "" {
		b.WriteString("**牌面要点**: ")
		b.WriteString(truncateRunes(desc, maxRunes))
		b.WriteString("\n\n")
	}
	if interpret != "" {
		b.WriteString("**牌义要点**: ")
		b.WriteString(truncateRunes(interpret, maxRunes))
		b.WriteString("\n\n")
	}
	return b.String(), nil
}

func truncateRunes(s string, max int) string {
	s = strings.TrimSpace(s)
	if max <= 0 {
		return s
	}
	runes := []rune(s)
	if len(runes) <= max {
		return s
	}
	return string(runes[:max]) + "…"
}

func buildCardMaterial(card request.CardInfo) (string, error) {
	detail, hasKB := tarotMap[card.Name]
	if !hasKB {
		if alt, ok := cardKeyAliases[card.Name]; ok {
			detail, hasKB = tarotMap[alt]
		}
	}
	if !hasKB {
		return "", fmt.Errorf("card %q not found in knowledge base", card.Name)
	}

	reversed := strings.EqualFold(card.Orientation, "reversed")
	state := "正位"
	keyword := detail.Upright
	interpretHeading := "Upright"
	interpretBody := ""
	if reversed {
		state = "逆位"
		keyword = detail.Reversed
		interpretHeading = "Reversed"
	}

	rag, hasRAG := lookupRAGChunk(card.Name)
	displayName := detail.Name
	if rag.Title != "" {
		displayName = rag.Title
	}
	if hasRAG {
		interpretBody = rag.Upright
		if reversed {
			interpretBody = rag.Reversed
		}
	}

	var b strings.Builder
	fmt.Fprintf(&b, "### %s · %s", displayName, state)
	if card.Meaning != "" {
		fmt.Fprintf(&b, " · 阵位：%s", card.Meaning)
	}
	b.WriteString("\n\n")

	fmt.Fprintf(&b, "**关键词（%s）**: %s\n\n", state, keyword)

	if hasRAG && rag.Description != "" {
		b.WriteString("**牌面描述**\n\n")
		b.WriteString(rag.Description)
		b.WriteString("\n\n")
	} else if detail.Description != "" {
		b.WriteString("**牌面描述**\n\n")
		b.WriteString(detail.Description)
		b.WriteString("\n\n")
	}

	if interpretBody != "" {
		fmt.Fprintf(&b, "**牌义解读（%s）**\n\n", interpretHeading)
		b.WriteString(interpretBody)
		b.WriteString("\n\n")
	}

	meta := []string{}
	if detail.Arcana != "" {
		meta = append(meta, "Arcana: "+detail.Arcana)
	}
	if detail.Element != "" {
		meta = append(meta, "Element: "+detail.Element)
	}
	if detail.Astrology != "" {
		meta = append(meta, "Astrology: "+detail.Astrology)
	}
	if detail.Numerology != "" {
		meta = append(meta, "Numerology: "+detail.Numerology)
	}
	if len(meta) > 0 {
		b.WriteString("*")
		b.WriteString(strings.Join(meta, " · "))
		b.WriteString("*\n")
	}

	return b.String(), nil
}
