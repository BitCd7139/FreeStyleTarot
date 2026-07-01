package api

import (
	"FreeStyleTarot/model/request"
	"FreeStyleTarot/service"
	"FreeStyleTarot/service/workflow/tarot"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func HandleClarify(c *gin.Context) {
	var req request.Clarify
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	qLen := len([]rune(req.Question))
	if qLen < 5 || qLen > 500 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "question length must be between 5 and 500 characters"})
		return
	}
	if len(req.Cards) > 15 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "too many cards"})
		return
	}
	for _, card := range req.Cards {
		if len([]rune(card.Meaning)) > 50 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "card meaning too long"})
			return
		}
	}

	if err := service.InitLlm(); err != nil {
		zap.S().Errorw("Failed to initialize LLM client", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "LLM unavailable"})
		return
	}

	ctx := c.Request.Context()
	if req.CustomAPI != nil {
		var err error
		ctx, err = withCustomAPI(ctx, req.CustomAPI)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "自定义 API 配置无效: " + err.Error()})
			return
		}
	}

	resp, err := tarot.RunIntentClarify(ctx, req)
	if err != nil {
		zap.S().Warnw("intent clarifier failed", "error", err)
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "意图分析暂时不可用，请稍后重试或直接开始占卜"})
		return
	}
	c.JSON(http.StatusOK, resp)
}
