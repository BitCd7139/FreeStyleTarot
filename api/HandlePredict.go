package api

import (
	"FreeStyleTarot/model/request"
	"FreeStyleTarot/model/response"
	"FreeStyleTarot/service"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func HandlePredict(c *gin.Context) {
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

	// 2. 调用 DeepSeek API
	llmResult, err := service.CallDeepSeek(c.Request.Context(), systemPrompt, userPrompt)
	if err != nil {
		zap.S().Errorw("DeepSeek call failed", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "AI service unavailable"})
		return
	}

	// 3. 组装最终返回给前端的内容
	finalAnswer := llmResult + "\n\n" + "--- \n" + "### 喜欢的话就来Github点个Star吧！项目链接：https://github.com/BitCd7139/FreeStyleTarot/tree/main"

	resp := &response.Predict{
		Answer: finalAnswer,
		Code:   200,
	}
	c.JSON(http.StatusOK, resp)
}
