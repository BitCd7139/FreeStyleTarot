package service

import (
	"fmt"
	"os"
	"sync"

	"github.com/tmc/langchaingo/llms/openai"
	"go.uber.org/zap"
)

// 依然使用你原来的全局变量名
var (
	apiKey      string
	apiBaseUrl  string
	apiModel    string
	apiInitErr  error
	apiInitOnce sync.Once

	globalLlm   *openai.LLM
	llmInitOnce sync.Once
)

func initApi() error {
	apiInitOnce.Do(func() {
		apiKey = os.Getenv("DEEPSEEK_API_KEY")
		apiBaseUrl = os.Getenv("DEEPSEEK_API_BASE")
		apiModel = os.Getenv("DEEPSEEK_API_MODEL")

		if apiKey == "" {
			apiInitErr = fmt.Errorf("必需的环境变量 DEEPSEEK_API_KEY 未设置")
			return
		}
		if apiBaseUrl == "" {
			apiInitErr = fmt.Errorf("必需的环境变量 DEEPSEEK_API_BASE 未设置")
			return
		}
		if apiModel == "" {
			apiInitErr = fmt.Errorf("必需的环境变量 DEEPSEEK_API_MODEL 未设置")
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
