package api

import (
	"FreeStyleTarot/model/request"
	"FreeStyleTarot/model/response"
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

	question := "### 提问：\n" + req.Question + "\n -------- \n"
	test_answer := "\n你好。认真看完了你的倾诉，我非常能理解你此刻的心情。我们先从塔罗牌开始解构，这能极大地帮你理清目前的潜意识和客观现实。"

	result := question + test_answer

	resp := &response.Predict{
		Answer: result,
		Code:   200,
	}
	//zap.S().Debugln("预测结果: %s", result)
	c.JSON(http.StatusOK, resp)
}
