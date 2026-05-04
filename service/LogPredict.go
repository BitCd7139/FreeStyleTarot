package service

import (
	"FreeStyleTarot/model/request"
	"fmt"

	"go.uber.org/zap"
)

func LogPredict(predict request.Predict) {
	zap.S().Infow("Predict input:", "model", predict.Model, "question", predict.Question, "Cards", len(predict.Cards))
	for _, card := range predict.Cards {
		zap.S().Infow("Card info:", "name", card.Name, "orientation", card.Orientation, "meaning", card.Meaning)
	}
	fmt.Print("\n")
}
