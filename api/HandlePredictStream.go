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

	llmInitErr := service.InitLlm()
	if llmInitErr != nil {
		panic("Failed to initialize LLM client: " + llmInitErr.Error())
		return
	}

	// 组装提示词
	systemPrompt, userPrompt, err := service.InputsAssembler(req)
	if err != nil {
		zap.S().Errorw("Failed to assemble card prompt", "card", req, "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Error"})
		return
	}

	// 2. 设置 SSE 响应头
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("X-Accel-Buffering", "no")

	clientContext := c.Request.Context()

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

		headerText := "## 🤔 提问\n" + req.Question + "\n\n---\n"
		if err := writeChunk(headerText); err != nil {
			zap.S().Errorw("Failed to write header", "error", err)
			return false
		}

		err := service.CallDeepSeekStream(clientContext, systemPrompt, userPrompt, func(chunk string) error {
			// 检查客户端是否已断开，避免浪费后续生成的 Token
			select {
			case <-clientContext.Done():
				return clientContext.Err()
			default:
				_, err := fmt.Fprintf(w, "data: %s\n\n", service.Trans2json(chunk))
				if err != nil {
					zap.S().Errorw("Failed to write to response stream", "error", err)
					return err
				}
				// 强刷缓冲区
				if flusher, ok := w.(http.Flusher); ok {
					flusher.Flush()
				}
				return nil
			}
		})

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
	zap.S().Infow("Predict input:", "question", req.Question)
}
