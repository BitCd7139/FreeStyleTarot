package api



import (

	"FreeStyleTarot/api/middleware"

	"FreeStyleTarot/model/request"
	"context"
	"errors"
	"unicode/utf8"

	"FreeStyleTarot/service"

	"FreeStyleTarot/service/workflow"

	"FreeStyleTarot/service/workflow/tarot"

	"fmt"

	"io"

	"net/http"

	"strings"



	"github.com/gin-gonic/gin"

	"go.uber.org/zap"

)



func HandlePredictStream(c *gin.Context) {

	var req request.Predict

	if err := c.ShouldBindJSON(&req); err != nil {

		zap.S().Errorw("Invalid request payload", "error", err)

		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})

		return

	}



	// 二次校验输入

	if utf8.RuneCountInString(req.Question) > 1500 {

		c.JSON(http.StatusBadRequest, gin.H{"error": "Input too long"})

		return

	}

	for _, card := range req.Cards {

		if len(card.Meaning) > 50 {

			c.JSON(http.StatusBadRequest, gin.H{"error": "Input too long"})

			return

		}

	}

	if len(req.Cards) > 15 {

		c.JSON(http.StatusBadRequest, gin.H{"error": "Input too long"})

		return

	}



	userID, ok := middleware.GetUserID(c)
	// 自定义 API 模式下不记录占卜时间，避免影响免费用户配额
	if ok && req.CustomAPI == nil {
		if err := authSvc.RecordPredict(c.Request.Context(), userID); err != nil {
			zap.S().Errorw("更新占卜时间失败", "error", err)
		}
	}



	llmInitErr := service.InitLlm()
	if llmInitErr != nil {
		panic("Failed to initialize LLM client: " + llmInitErr.Error())
		return
	}

	clientContext := c.Request.Context()
	if req.CustomAPI != nil {
		var err error
		clientContext, err = withCustomAPI(clientContext, req.CustomAPI)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "自定义 API 配置无效: " + err.Error()})
			return
		}
	}

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")

	var answerBuilder strings.Builder



	c.Stream(func(w io.Writer) bool {

		writeEvent := func(event workflow.Event) error {

			select {

			case <-clientContext.Done():

				return clientContext.Err()

			default:

				if text, ok := workflow.AppendAnswerText(event); ok {

					answerBuilder.WriteString(text)

				}

				return workflow.WriteEvent(w, event)

			}

		}



		introText := "## 🤔 用户提问\n" + req.Question + tarot.FormatClarificationIntro(req) + "\n\n---\n"

		if err := writeEvent(workflow.IntroEvent(introText)); err != nil {

			zap.S().Errorw("Failed to write intro", "error", err)

			return false

		}



		if err := tarot.RunReading(clientContext, req, writeEvent); err != nil {

			if errors.Is(err, context.Canceled) {
				zap.S().Infow("Predict stream closed by client")
				return false
			}

			zap.S().Errorw("Streaming error occurred", "error", err)

			fmt.Fprintf(w, "event: error\ndata: {\"message\": \"Internal server error during generation\"}\n\n")

			return false

		}

		outroText := "\n\n---\n### 喜欢的话就来Github点个Star吧！https://github.com/BitCd7139/FreeStyleTarot " + "\n ### 如果可以的话也来B站点赞投个币吧！https://www.bilibili.com/video/BV1gSReBJELE/"

		if err := writeEvent(workflow.OutroEvent(outroText)); err != nil {

			zap.S().Errorw("Failed to write outro", "error", err)

			return false

		}



		if ok {

			if err := historySvc.SaveFromPredict(clientContext, userID, req, tarot.HistoryAnswer(req, answerBuilder.String())); err != nil {

				zap.S().Errorw("保存提问历史失败", "error", err)

			}

		}



		fmt.Fprintf(w, "event: close\ndata: [DONE]\n\n")

		zap.S().Infow("Stream session closed cleanly")

		return false

	})

	service.LogPredict(req)

}

