package service

import (
	"FreeStyleTarot/model/prompt"
	"FreeStyleTarot/model/request"
	"FreeStyleTarot/storage"
	"encoding/xml"
	"errors"
	"strings"
	"sync"

	"go.uber.org/zap"
)

var (
	// 全局缓存 map
	backgroundPrompts map[string][]byte
	// 确保初始化逻辑只运行一次
	backgroundPromptOnce sync.Once
)

const (
	prefix = "background_prompt_"
	suffix = ".md"
)

func InputsAssembler(predict request.Predict) (systemMsg string, userMsg string, err error) {
	loadKnowledgeBase()

	// 1. 系统角色/背景 (System Prompt)
	systemMsg = string(getBackgroundPrompt(predict.Model))

	// 2. 用户输入内容 (User Prompt)
	var userContent strings.Builder
	for _, card := range predict.Cards {
		p, err := cardPromptAssembler(card)
		if err != nil {
			zap.S().Errorw("Failed to assemble card prompt", "card", card.Name, "error", err)
			continue
		}
		userContent.WriteString(p + "\n\n")
	}

	userContent.WriteString("\n## Question:\n" + predict.Question + "\n")

	return systemMsg, userContent.String(), nil
}

func cardPromptAssembler(card request.CardInfo) (string, error) {

	detail, exists := tarotMap[card.Name]
	if !exists {
		zap.S().Errorw("card %s not found in knowledge base", card.Name)
		return "", errors.New("card not found in knowledge base")
	}

	// 3. 根据 Orientation 决定关键词
	keyword := detail.Upright
	state := "正位"
	if strings.ToLower(card.Orientation) == "reversed" {
		keyword = detail.Reversed
		state = "逆位"
	}

	// 4. 组装输出结构体
	cardxml := prompt.CardPrompt{
		Name:        detail.Name,
		Meaning:     card.Meaning,
		State:       state,
		Description: detail.Description,
		KeyWord:     keyword,
		Arcana:      detail.Arcana,
		Element:     detail.Element,
		Numerology:  detail.Numerology,
		Astrology:   detail.Astrology,
	}

	output, err := xml.MarshalIndent(cardxml, "", "  ")
	if err != nil {
		return "", err
	}

	return string(output), nil
}

func initBackgroundPrompts() {
	backgroundPromptOnce.Do(func() {
		backgroundPrompts = make(map[string][]byte)

		entries, err := storage.Assets.ReadDir(".")
		if err != nil {
			zap.S().Errorf("CRITICAL: Failed to read assets directory: %v", err)
			return
		}

		zap.S().Debug("Starting to scan background prompts...")

		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}

			filename := entry.Name()

			// 3. 校验文件名格式：必须以 background_prompt_ 开头，以 .md 结尾
			if strings.HasPrefix(filename, prefix) && strings.HasSuffix(filename, suffix) {

				// 4. 提取 name 部分
				// 比如 "background_prompt_customer_service.md" -> "customer_service"
				name := strings.TrimPrefix(filename, prefix)
				name = strings.TrimSuffix(name, suffix)

				// 5. 读取文件内容
				data, err := storage.Assets.ReadFile(filename)
				if err != nil {
					zap.S().Errorf("Failed to read file %s: %v", filename, err)
					continue
				}

				// 6. 存入 Map
				backgroundPrompts[name] = data
				zap.S().Infow("Loaded background prompt", "name", name, "file", filename)
			}
		}

		zap.S().Infof("Initialization complete. Loaded %d prompts.", len(backgroundPrompts))
	})
}

func getBackgroundPrompt(name string) []byte {
	initBackgroundPrompts()

	if data, ok := backgroundPrompts[name]; ok {
		return data
	}

	zap.S().Warnw("Prompt not found", "name", name)
	return nil
}

func targetPrompt() string {
	return "### 喜欢的话就来Github点个Star吧！项目链接：https://github.com/BitCd7139/FreeStyleTarot/tree/main"
}
