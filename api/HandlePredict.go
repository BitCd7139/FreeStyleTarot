package api

import (
	"FreeStyleTarot/model/request"
	"FreeStyleTarot/model/response"
	"FreeStyleTarot/service"
	"html"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func HandlePredict(c *gin.Context) {
	var req request.Predict

	if err := c.ShouldBindJSON(&req); err != nil {
		zap.S().Errorw("请求参数绑定失败: %v", err)
		zap.S().Errorw("请求参数: %v", req)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// result := "test result"
	// api_key = os.getenv("OPENAI_API_KEY")
	// api_link = os.getnvar("OPENAI_API_LINK")

	result, err := service.InputsAssembler(req)

	if err != nil {
		zap.S().Error(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if result == "" {
		zap.S().Error("Invalid input")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	//result := question + answer
	result = html.EscapeString(result)

	resp := &response.Predict{
		Answer: result,
		Code:   200,
	}
	//zap.S().Debugln("预测结果: %s", result)
	c.JSON(http.StatusOK, resp)
}
