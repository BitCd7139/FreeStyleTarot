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

func InputsAssembler(predict request.Predict) (string, error) {
	loadKnowledgeBase()
	loadBackgroundPrompt()

	output := string(backgroundPrompt) + "\n\n"

	for _, card := range predict.Cards {
		prompt, err := cardPromptAssembler(card)
		if err != nil {
			zap.S().Errorw("Failed to assemble card prompt for card %s: %v", card.Name, err)
			continue
		}
		output += prompt + "\n\n"
	}

	output += targetPrompt()
	return output, nil
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
		zap.S().Debugw("BackgroundPrompt loaded", "content", string(backgroundPrompt))
	}
}

func targetPrompt() string {
	return "#请根据以上牌面信息，结合提问者的问题和角色限制条件，进行塔罗牌解读。"
}
