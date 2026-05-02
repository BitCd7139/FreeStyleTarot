package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/tmc/langchaingo/llms"
	"go.uber.org/zap"
)

func CallDeepSeek(ctx context.Context, systemPrompt string, userPrompt string) (string, error) {

	// 2. 构建消息组合
	content := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, systemPrompt),
		llms.TextParts(llms.ChatMessageTypeHuman, userPrompt),
	}

	// 3. 调用生成
	resp, err := globalLlm.GenerateContent(ctx, content,
		llms.WithTemperature(0.7),
		llms.WithMaxTokens(2048),
	)
	if err != nil {
		return "", err
	}

	if len(resp.Choices) > 0 {
		return resp.Choices[0].Content, nil
	}

	return "", fmt.Errorf("no response from deepseek")
}

func CallDeepSeekStream(ctx context.Context, systemPrompt string, userPrompt string, onChunk func(chunk string) error) error {
	if globalLlm == nil {
		return fmt.Errorf("LLM client not initialized")
	}

	content := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, systemPrompt),
		llms.TextParts(llms.ChatMessageTypeHuman, userPrompt),
	}

	ctx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	defer cancel()

	// 调用 GenerateContent，并传入 WithStreamingFunc
	_, err := globalLlm.GenerateContent(ctx, content,
		llms.WithMaxTokens(2048),
		llms.WithTemperature(0.7),
		// 关键点：每当收到一个 token，这个函数就会被触发
		llms.WithStreamingFunc(func(ctx context.Context, chunk []byte) error {
			select {
			case <-ctx.Done():
				zap.S().Warn("Stream processing interrupted by context cancellation")
				return ctx.Err()
			default:
				return onChunk(string(chunk))
			}
		}),
	)

	if err != nil {
		if errors.Is(err, context.Canceled) {
			zap.S().Warn("DeepSeek request was canceled by the client")
		} else {
			zap.S().Errorw("DeepSeek API error", "error", err)
		}
		return err
	}

	zap.S().Debugw("LLM streaming request completed successfully")
	return nil
}

func Trans2json(input string) []byte {
	data := map[string]interface{}{
		"content": input,
	}
	jsonData, _ := json.Marshal(data)
	return jsonData
}
