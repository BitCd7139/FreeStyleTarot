package api

import (
	"FreeStyleTarot/model/request"
	"FreeStyleTarot/service"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func HandlePrompt(c *gin.Context) {
	var req request.Predict
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 1. 组装提示词
	systemPrompt, userPrompt, err := service.InputsAssembler(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to assemble prompt"})
		return
	}

	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")

	// 3. 开始流式写入
	c.Stream(func(w io.Writer) bool {
		writeChunk := func(content string) error {
			// 转义为 JSON 字符串
			_, err := fmt.Fprintf(w, "data: %s\n\n", service.Trans2json(content))
			if err != nil {
				return err
			}
			if flusher, ok := w.(http.Flusher); ok {
				flusher.Flush()
			}
			return nil
		}

		headerText := systemPrompt + userPrompt
		if err := writeChunk(headerText); err != nil {
			zap.S().Errorw("Failed to write header", "error", err)
			return false
		}

		if err != nil {
			zap.S().Errorw("Streaming error occurred", "error", err)
			fmt.Fprintf(w, "event: error\ndata: {\"message\": \"Internal server error during generation\"}\n\n")
			return false
		}

		footerText := "\n\n---\n### 喜欢的话就来Github点个Star吧！项目链接：https://github.com/BitCd7139/FreeStyleTarot"
		if err := writeChunk(footerText); err != nil {
			zap.S().Errorw("Failed to write footer", "error", err)
			return false
		}

		// 发送结束标识
		fmt.Fprintf(w, "event: close\ndata: [DONE]\n\n")
		zap.S().Infow("Stream session closed cleanly")
		return false
	})
	zap.S().Infow("Predict input:", "model", req.Model, "question", req.Question)
}
