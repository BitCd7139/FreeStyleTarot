package api

import (
	"FreeStyleTarot/api/middleware"
	"FreeStyleTarot/model/request"
	"FreeStyleTarot/service"
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
	if len(req.Question) > 2000 {
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
	if ok {
		if err := authSvc.RecordPredict(c.Request.Context(), userID); err != nil {
			zap.S().Errorw("更新占卜时间失败", "error", err)
		}
	}

	llmInitErr := service.InitLlm()
	if llmInitErr != nil {
		panic("Failed to initialize LLM client: " + llmInitErr.Error())
		return
	}

	// 2. 设置 SSE 响应头
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")

	clientContext := c.Request.Context()
	var answerBuilder strings.Builder

	// 3. 开始流式写入
	c.Stream(func(w io.Writer) bool {
		writeEvent := func(event service.StreamEvent) error {
			select {
			case <-clientContext.Done():
				return clientContext.Err()
			default:
				if event.Type == service.StreamEventTypeContent && event.Content != "" {
					answerBuilder.WriteString(event.Content)
				}
				if err := service.WriteStreamEvent(w, event); err != nil {
					return err
				}
				return nil
			}
		}

		headerText := "## 🤔 提问\n" + req.Question + "\n\n---\n"
		if err := writeEvent(service.ContentEvent(headerText)); err != nil {
			zap.S().Errorw("Failed to write header", "error", err)
			return false
		}

		materials, err := service.BuildSpreadMaterialsParallel(req.Cards)
		if err != nil {
			zap.S().Errorw("Failed to build spread materials", "error", err)
			fmt.Fprintf(w, "event: error\ndata: {\"message\": \"Failed to build spread materials\"}\n\n")
			return false
		}

		err = service.RunTarotReading(clientContext, req, materials, writeEvent)

		if err != nil {
			zap.S().Errorw("Streaming error occurred", "error", err)
			fmt.Fprintf(w, "event: error\ndata: {\"message\": \"Internal server error during generation\"}\n\n")
			return false
		}

		footerText := "\n\n---\n### 喜欢的话就来Github点个Star吧！项目链接：https://github.com/BitCd7139/FreeStyleTarot " + "\n ### 如果可以的话也来 B站点赞投个币吧！视频链接：https://www.bilibili.com/video/BV1gSReBJELE/"
		if err := writeEvent(service.ContentEvent(footerText)); err != nil {
			zap.S().Errorw("Failed to write footer", "error", err)
			return false
		}

		if ok {
			if err := historySvc.SaveFromPredict(clientContext, userID, req, answerBuilder.String()); err != nil {
				zap.S().Errorw("保存提问历史失败", "error", err)
			}
		}

		// 发送结束标识
		fmt.Fprintf(w, "event: close\ndata: [DONE]\n\n")
		zap.S().Infow("Stream session closed cleanly")
		return false
	})
	service.LogPredict(req)
}
