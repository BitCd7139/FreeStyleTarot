package service

import (
	"context"
	"fmt"
	"strings"

	"FreeStyleTarot/model/request"

	"github.com/tmc/langchaingo/llms"
	"go.uber.org/zap"
)

// RunTarotReading runs Agent A (analysis) then streams Agent B (persona advice).
func RunTarotReading(ctx context.Context, predict request.Predict, onChunk func(chunk string) error) error {
	materials, err := BuildSpreadMaterials(predict.Cards)
	if err != nil {
		return fmt.Errorf("build spread materials: %w", err)
	}
	brief, err := BuildSpreadBriefMaterials(predict.Cards)
	if err != nil {
		return fmt.Errorf("build spread brief materials: %w", err)
	}

	systemPrompt := SharedSystemPrompt(predict.Model)
	if systemPrompt == "" {
		return fmt.Errorf("shared system prompt not loaded")
	}

	analysis, err := callSpreadAnalyst(ctx, systemPrompt, AgentAUserPrompt(predict, materials))
	if err != nil {
		return fmt.Errorf("agent A: %w", err)
	}

	bUser := AgentBUserPrompt(predict, brief, analysis)
	zap.S().Debugw("Streaming Agent B", "model", predict.Model)
	return CallDeepSeekStream(ctx, systemPrompt, bUser, onChunk,
		agentPersonaAdvisor.MaxTokens, agentPersonaAdvisor.Temperature)
}

// RunTarotReadingSync runs Agent A then Agent B without streaming.
func RunTarotReadingSync(ctx context.Context, predict request.Predict) (string, error) {
	materials, err := BuildSpreadMaterials(predict.Cards)
	if err != nil {
		return "", fmt.Errorf("build spread materials: %w", err)
	}
	brief, err := BuildSpreadBriefMaterials(predict.Cards)
	if err != nil {
		return "", fmt.Errorf("build spread brief materials: %w", err)
	}

	systemPrompt := SharedSystemPrompt(predict.Model)
	if systemPrompt == "" {
		return "", fmt.Errorf("shared system prompt not loaded")
	}

	analysis, err := callSpreadAnalyst(ctx, systemPrompt, AgentAUserPrompt(predict, materials))
	if err != nil {
		return "", fmt.Errorf("agent A: %w", err)
	}

	out, err := callAgent(ctx, systemPrompt, AgentBUserPrompt(predict, brief, analysis),
		agentPersonaAdvisor.Temperature, agentPersonaAdvisor.MaxTokens)
	if err != nil {
		return "", fmt.Errorf("agent B: %w", err)
	}
	return out, nil
}

func callSpreadAnalyst(ctx context.Context, systemPrompt, userPrompt string) (string, error) {
	zap.S().Debug("Running Agent A (spread_analyst)")
	analysis, err := callAgent(ctx, systemPrompt, userPrompt,
		agentSpreadAnalyst.Temperature, agentSpreadAnalyst.MaxTokens)
	if err != nil {
		return "", err
	}
	if strings.TrimSpace(analysis) == "" {
		return "", fmt.Errorf("empty analysis")
	}
	return analysis, nil
}

func callAgent(ctx context.Context, systemPrompt, userPrompt string, temperature float64, maxTokens int) (string, error) {
	if globalLlm == nil {
		return "", fmt.Errorf("LLM client not initialized")
	}

	content := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, systemPrompt),
		llms.TextParts(llms.ChatMessageTypeHuman, userPrompt),
	}

	resp, err := globalLlm.GenerateContent(ctx, content,
		llms.WithTemperature(temperature),
		llms.WithMaxTokens(maxTokens),
	)
	if err != nil {
		return "", err
	}
	if len(resp.Choices) > 0 {
		return resp.Choices[0].Content, nil
	}
	return "", fmt.Errorf("no response from LLM")
}
