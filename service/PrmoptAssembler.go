package service

import (
	"FreeStyleTarot/model/prmopt"
	"FreeStyleTarot/model/request"
	"encoding/xml"
	"errors"
	"strings"

	"go.uber.org/zap"
)

func InputsAssembler(predict request.Predict) (string, error) {
	output := backgroundPrmopt()

	for _, card := range predict.Cards {
		prompt, err := cardPromptAssembler(card)
		if err != nil {
			zap.S().Errorw("Failed to assemble card prompt for card %s: %v", card.Name, err)
			continue
		}
		output += prompt + "\n\n"
	}

	output += targetPrmopt()
	return output, nil
}

func cardPromptAssembler(card request.CardInfo) (string, error) {
	loadKnowledgeBase()

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
	prompt := prmopt.CardPrmopt{
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

	output, err := xml.MarshalIndent(prompt, "", "  ")
	if err != nil {
		return "", err
	}

	return string(output), nil
}

func backgroundPrmopt() string {
	return
}

func targetPrmopt() string {
	return "### 目标：\n" +
		"请根据以上牌面信息，结合提问者的问题和角色限制条件，进行塔罗牌解读。"
}
