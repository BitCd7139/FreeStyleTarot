package service

import (
	"FreeStyleTarot/model/prompt"
	"FreeStyleTarot/model/request"
	"FreeStyleTarot/storage"
	"encoding/xml"
	"errors"
	"fmt"
	"strings"
	"sync"

	"go.uber.org/zap"
)

var backgroundPrompt []byte
var backgroundPromptOnce sync.Once

func InputsAssembler(predict request.Predict) (systemMsg string, userMsg string, err error) {
	loadKnowledgeBase()
	loadBackgroundPrompt()

	// 1. 系统角色/背景 (System Prompt)
	systemMsg = string(backgroundPrompt)

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

func loadBackgroundPrompt() {
	entries, _ := storage.Assets.ReadDir(".")
	for _, entry := range entries {
		fmt.Println("Found file:", entry.Name())
	}

	backgroundPromptOnce.Do(func() {
		data, err := storage.Assets.ReadFile("background_prompt.md")
		if err != nil {
			zap.S().Errorf("CRITICAL: Failed to read file: %v", err)
			return
		}
		backgroundPrompt = data
	})

	if len(backgroundPrompt) == 0 {
		zap.S().Warn("BackgroundPrompt is empty!")
	} else {
		//zap.S().Debugw("BackgroundPrompt loaded", "content", string(backgroundPrompt))
	}
}

func targetPrompt() string {
	return "### 喜欢的话就来Github点个Star吧！项目链接：https://github.com/BitCd7139/FreeStyleTarot/tree/main"
}
