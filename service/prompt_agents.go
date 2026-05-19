package service

import (
	"FreeStyleTarot/model/request"
	"FreeStyleTarot/storage"
	"fmt"
	"strings"
	"sync"

	"go.uber.org/zap"
)

const (
	roleSpreadAnalyst  = "spread_analyst"
	rolePersonaAdvisor = "persona_advisor"
)

var (
	sharedSystemOnce sync.Once
	sharedSystemBase []byte
)

func loadSharedSystemBase() {
	sharedSystemOnce.Do(func() {
		data, err := storage.Assets.ReadFile("tarot_shared_system.md")
		if err != nil {
			zap.S().Errorw("Failed to load tarot_shared_system.md", "error", err)
			return
		}
		sharedSystemBase = data
	})
}

// SharedSystemPrompt is byte-identical for Agent A and Agent B when model is fixed (DeepSeek prefix cache).
func SharedSystemPrompt(model string) string {
	loadSharedSystemBase()
	initBackgroundPrompts()

	var b strings.Builder
	b.Write(sharedSystemBase)
	persona := getBackgroundPrompt(model)
	if len(persona) > 0 {
		b.Write(persona)
	}
	return b.String()
}

// activeAgentCue is a fixed short prefix; role details live only in SharedSystemPrompt.
func activeAgentCue(role string) string {
	return "[ACTIVE AGENT: " + role + "]\n\n"
}

func AgentAUserPrompt(predict request.Predict, materials string) string {
	var b strings.Builder
	b.WriteString(activeAgentCue(roleSpreadAnalyst))
	b.WriteString(materials)
	b.WriteString("\n\n---\n\n## Question\n\n")
	b.WriteString(predict.Question)
	b.WriteString("\n")
	return b.String()
}

func AgentBUserPrompt(predict request.Predict, briefMaterials, spreadAnalysis string) string {
	var b strings.Builder
	b.WriteString(activeAgentCue(rolePersonaAdvisor))
	b.WriteString("## User question\n\n")
	b.WriteString(predict.Question)
	b.WriteString("\n\n---\n\n")
	if briefMaterials != "" {
		b.WriteString(briefMaterials)
		b.WriteString("\n\n---\n\n")
	}
	b.WriteString("## Spread analysis (from spread_analyst)\n\n")
	b.WriteString(spreadAnalysis)
	b.WriteString("\n")
	return b.String()
}

func DebugPromptBundle(predict request.Predict) (string, error) {
	materials, err := BuildSpreadMaterials(predict.Cards)
	if err != nil {
		return "", err
	}
	sys := SharedSystemPrompt(predict.Model)
	aUser := AgentAUserPrompt(predict, materials)
	brief, _ := BuildSpreadBriefMaterials(predict.Cards)
	bUser := AgentBUserPrompt(predict, brief, "« Agent A output »")

	var b strings.Builder
	fmt.Fprintf(&b, "=== Shared system (both agents) ===\n%s\n\n", sys)
	fmt.Fprintf(&b, "=== Agent A — user ===\n%s\n\n", aUser)
	fmt.Fprintf(&b, "=== Agent B — user ===\n%s\n", bUser)
	return b.String(), nil
}
