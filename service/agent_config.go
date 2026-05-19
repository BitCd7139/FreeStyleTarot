package service

import (
	"FreeStyleTarot/config"

	"go.uber.org/zap"
)

type agentRuntimeConfig struct {
	MaxTokens   int
	Temperature float64
}

var (
	agentSpreadAnalyst  agentRuntimeConfig
	agentPersonaAdvisor agentRuntimeConfig
	briefMaterialMaxRunes int
)

// InitAgentConfig loads agent LLM limits from config and warms tarot indexes.
func InitAgentConfig() {
	cfg := config.GlobalConfig.Agents

	agentSpreadAnalyst = agentRuntimeConfig{
		MaxTokens:   orDefault(cfg.SpreadAnalyst.MaxTokens, 4096),
		Temperature: orDefaultFloat(cfg.SpreadAnalyst.Temperature, 0.6),
	}
	agentPersonaAdvisor = agentRuntimeConfig{
		MaxTokens:   orDefault(cfg.PersonaAdvisor.MaxTokens, 8192),
		Temperature: orDefaultFloat(cfg.PersonaAdvisor.Temperature, 0.7),
	}
	briefMaterialMaxRunes = orDefault(cfg.BriefMaterialMaxRunes, 500)

	loadKnowledgeBase()
	loadTarotRAG()
	initBackgroundPrompts()

	zap.S().Infow("Agent config initialized",
		"spread_analyst_max_tokens", agentSpreadAnalyst.MaxTokens,
		"persona_advisor_max_tokens", agentPersonaAdvisor.MaxTokens,
		"brief_material_max_runes", briefMaterialMaxRunes,
	)
}

func orDefault(v, def int) int {
	if v <= 0 {
		return def
	}
	return v
}

func orDefaultFloat(v, def float64) float64 {
	if v <= 0 {
		return def
	}
	return v
}
