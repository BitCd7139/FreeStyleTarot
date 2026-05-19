package service

import (
	"FreeStyleTarot/storage"
	"encoding/json"
	"fmt"
	"sync"

	"go.uber.org/zap"
)

type TarotDetail struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Upright     string `json:"upright"`
	Reversed    string `json:"reversed"`
	Arcana      string `json:"arcana"`
	Element     string `json:"element"`
	Numerology  string `json:"numerology"`
	Astrology   string `json:"astrology"`
}

var (
	tarotMap     map[string]TarotDetail
	knowbaseOnce sync.Once
)

func loadKnowledgeBase() {
	knowbaseOnce.Do(func() {
		tarotMap = make(map[string]TarotDetail)

		data, err := storage.Assets.ReadFile("tarot_knowledge_base.json")
		if err != nil {
			fmt.Printf("Error reading knowledge base: %v\n", err)
			return
		}
		if err := json.Unmarshal(data, &tarotMap); err != nil {
			fmt.Printf("Error unmarshaling knowledge base: %v\n", err)
		}
		zap.S().Infow("Loaded knowledge base", "data", len(tarotMap))

	})
}
