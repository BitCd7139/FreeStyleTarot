package api

import (
	"net/http"

	"FreeStyleTarot/api/middleware"
	"FreeStyleTarot/service/history"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var historySvc = history.NewService()

// HandleGetPredictHistory GET /auth/predict-history
func HandleGetPredictHistory(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未登录"})
		return
	}

	result, err := historySvc.ListForUser(c.Request.Context(), userID)
	if err != nil {
		zap.S().Errorw("获取提问历史失败", "error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取提问历史失败"})
		return
	}

	c.JSON(http.StatusOK, result)
}
