package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/tmc/langchaingo/llms"
	"github.com/tmc/langchaingo/llms/openai"
	"go.uber.org/zap"
)

var (
	apiKey      string
	apiBaseUrl  string
	apiModel    string
	apiInitOnce sync.Once
	apiInitErr  error

	globalLlm   llms.Model
	llmInitOnce sync.Once
)

func initApi() error {
	apiInitOnce.Do(func() {
		apiKey = os.Getenv("DEEPSEEK_API_KEY")
		if apiKey == "" {
			apiInitErr = fmt.Errorf("DEEPSEEK_API_KEY not set")
			return
		}
		apiBaseUrl = os.Getenv("DEEPSEEK_API_BASE")
		if apiBaseUrl == "" {
			apiInitErr = fmt.Errorf("DEEPSEEK_API_BASE not set")
			return
		}
		apiModel = os.Getenv("DEEPSEEK_API_MODEL")
		if apiModel == "" {
			apiInitErr = fmt.Errorf("DEEPSEEK_API_MODEL not set")
			return
		}
	})
	return apiInitErr
}

func InitLlm() error {
	err := initApi()
	if err != nil {
		return err
	}

	llmInitOnce.Do(func() {
		globalLlm, err = openai.New(
			openai.WithToken(apiKey),
			openai.WithBaseURL(apiBaseUrl),
			openai.WithModel(apiModel),
		)
		if err != nil {
			zap.S().Errorw("Failed to initialize DeepSeek client", "error", err)
			return
		}
	})
	return err
}

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
