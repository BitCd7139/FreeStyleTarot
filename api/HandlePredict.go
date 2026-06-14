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

	if _, err := service.BuildSpreadMaterials(req.Cards); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	llmResult, err := service.RunTarotReadingSync(c.Request.Context(), req)
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
